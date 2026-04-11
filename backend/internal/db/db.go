package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "embed"
)

//go:embed sql/init.sql
var initSQL string

func InitDB(db *sql.DB) error {
	_, err := db.Exec(initSQL)
	return err
}

func GetPostgreSQL_db(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	if err := InitDB(db); err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return db, nil
}
