package seeder

import (
	"database/sql"
	"io/fs"
	"log/slog"
)

type Seeder struct {
	DB   *sql.DB
	fsys fs.FS
	log  *slog.Logger
}

func New(DB *sql.DB, fsys fs.FS, logger *slog.Logger) *Seeder {
	return &Seeder{
		DB:   DB,
		fsys: fsys,
		log:  logger,
	}
}

func (s *Seeder) Seed() error {
	entries, err := fs.ReadDir(s.fsys, ".")
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, e := range entries {
		data, err := fs.ReadFile(s.fsys, e.Name())
		if err != nil {
			return err
		}

		s.log.Info("applying seed", "name", e.Name())

		_, err = tx.Exec(string(data))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
