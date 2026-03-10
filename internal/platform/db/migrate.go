package db

import (
	"cmp"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"slices"
	"strings"

	"github.com/yancarlodev/workspaces-api/internal/platform/datastruct"
)

type migrator struct {
	DB   *sql.DB
	fsys fs.FS
	log  *slog.Logger
}

func NewMigrator(DB *sql.DB, fsys fs.FS, logger *slog.Logger) *migrator {
	return &migrator{
		fsys: fsys,
		DB:   DB,
		log:  logger,
	}
}

func (m *migrator) Migrate() error {
	if err := m.createTable(); err != nil {
		return err
	}

	executed, err := m.executedMigrations()
	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(m.fsys, ".")
	if err != nil {
		return err
	}

	pending := m.pendingMigrations(entries, executed)
	if len(pending) == 0 {
		m.log.Info("no pending migration found")
		return nil
	}

	m.log.Info("pending migration", "count", len(pending))

	return m.exec(pending)
}

func (m *migrator) createTable() error {
	_, err := m.DB.Exec(`CREATE TABLE IF NOT EXISTS migration (
		name TEXT PRIMARY KEY,
		checksum TEXT,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)

	return err
}

func (m *migrator) executedMigrations() (datastruct.Set[string], error) {
	rows, err := m.DB.Query("SELECT name FROM migration;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := datastruct.NewHashSet[string]()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		executed.Add(name)
	}

	return executed, nil
}

func (m *migrator) pendingMigrations(entries []fs.DirEntry, executed datastruct.Set[string]) []fs.DirEntry {
	upMigrationCount := len(entries) / 2
	pendingCount := upMigrationCount - executed.Size()

	pending := make([]fs.DirEntry, 0, pendingCount)
	for _, entry := range entries {
		isUpMigration := strings.HasSuffix(entry.Name(), ".up.sql")

		if isUpMigration && !executed.Contains(entry.Name()) {
			pending = append(pending, entry)
		}
	}

	slices.SortFunc(pending, func(a, b fs.DirEntry) int {
		return cmp.Compare(a.Name(), b.Name())
	})

	return pending
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
		m.log.Info("applying migration", "name", entryName)

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
