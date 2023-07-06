package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://root:password@localhost:5433/go-chat?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (database *Database) Close() {
	database.DB.Close()
}

func (database *Database) GetDatabase() *sql.DB {
	return database.DB
}
