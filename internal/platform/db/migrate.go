package db

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/yancarlodev/workspaces-api/internal/platform/datastruct"
)

func Migrate(DB *sql.DB, fsys *embed.FS, sqlFilePath string) error {
	sub, entries, err := openSQLDir(fsys, sqlFilePath)
	if err != nil {
		return err
	}

	if err = createMigrationTable(DB); err != nil {
		return err
	}

	executedMigrations, err := getExecutedMigrations(DB)
	if err != nil {
		return err
	}

	pendingMigration := getPendingMigrations(entries, executedMigrations)

	if err = execMigrations(DB, sub, pendingMigration); err != nil {
		return err
	}

	return nil
}

func openSQLDir(fsys *embed.FS, sqlFilePath string) (fs.FS, []fs.DirEntry, error) {
	sub, err := fs.Sub(fsys, sqlFilePath)
	if err != nil {
		return nil, nil, err
	}

	entries, err := fs.ReadDir(sub, ".")
	if err != nil {
		return nil, nil, err
	}

	return sub, entries, nil
}

func createMigrationTable(DB *sql.DB) error {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS migration (
		name TEXT PRIMARY KEY,
		checksum TEXT,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)

	return err
}

func getExecutedMigrations(DB *sql.DB) (datastruct.Set[string], error) {
	rows, err := DB.Query("SELECT name FROM migration;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executedMigrations := datastruct.NewHashSet[string]()
	for rows.Next() {
		var migrationName string
		if err := rows.Scan(&migrationName); err != nil {
			return nil, err
		}

		executedMigrations.Add(migrationName)
	}

	return executedMigrations, nil
}

func getPendingMigrations(entries []fs.DirEntry, executedMigrations datastruct.Set[string]) []fs.DirEntry {
	upMigrationCount := len(entries) / 2
	pendingMigrationsCount := upMigrationCount - executedMigrations.Size()

	pendingMigration := make([]fs.DirEntry, 0, pendingMigrationsCount)
	for _, entry := range entries {
		isUpMigration := strings.HasSuffix(entry.Name(), "up.sql")

		if isUpMigration && !executedMigrations.Contains(entry.Name()) {
			pendingMigration = append(pendingMigration, entry)
		}
	}

	return pendingMigration
}

func execMigrations(DB *sql.DB, fsys fs.FS, entries []fs.DirEntry) error {
	stmt, err := DB.Prepare(`
		INSERT INTO migration(name, checksum)
		VALUES ($1, $2);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		entryName := entry.Name()

		data, err := fs.ReadFile(fsys, entryName)
		if err != nil {
			return err
		}

		hasher := sha256.New()
		_, err = hasher.Write(data)
		if err != nil {
			return err
		}

		checksum := hasher.Sum(nil)

		_, err = DB.Exec(string(data))
		if err != nil {
			return err
		}

		_, err = stmt.Exec(entryName, fmt.Sprintf("%x", checksum))
		if err != nil {
			return err
		}
	}

	return nil
}
