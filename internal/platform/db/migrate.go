package db

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/fs"
	"strings"

	"github.com/yancarlodev/workspaces-api/internal/platform/datastruct"
)

type migrator struct {
	DB   *sql.DB
	fsys fs.FS
}

func NewMigrator(DB *sql.DB, fsys fs.FS) *migrator {
	return &migrator{
		fsys: fsys,
		DB:   DB,
	}
}

func (m *migrator) Migrate() error {
	if err := m.createTable(); err != nil {
		return err
	}

	executedMigrations, err := m.getExecutedMigrations()
	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(m.fsys, ".")
	if err != nil {
		return err
	}

	pendingMigration := m.getPendingMigrations(entries, executedMigrations)

	return m.exec(pendingMigration)
}

func (m *migrator) createTable() error {
	_, err := m.DB.Exec(`CREATE TABLE IF NOT EXISTS migration (
		name TEXT PRIMARY KEY,
		checksum TEXT,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)

	return err
}

func (m *migrator) getExecutedMigrations() (datastruct.Set[string], error) {
	rows, err := m.DB.Query("SELECT name FROM migration;")
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

func (m *migrator) getPendingMigrations(entries []fs.DirEntry, executedMigrations datastruct.Set[string]) []fs.DirEntry {
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

func (m *migrator) exec(entries []fs.DirEntry) error {
	stmt, err := m.DB.Prepare(`
		INSERT INTO migration(name, checksum)
		VALUES ($1, $2);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		entryName := entry.Name()

		data, err := fs.ReadFile(m.fsys, entryName)
		if err != nil {
			return err
		}

		hasher := sha256.New()
		_, err = hasher.Write(data)
		if err != nil {
			return err
		}

		checksum := hasher.Sum(nil)

		_, err = m.DB.Exec(string(data))
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
