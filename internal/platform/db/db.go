package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Init(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %e", err)
	}

	return db, nil
}
