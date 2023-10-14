package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type DB struct {
	db *sql.DB
}

func NewDB() (*DB, error) {
	connectionString := "host=localhost port=5432 user=postgres password=postgres dbname=hotelreservation sslmode=disable"

	connectionPool, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := connectionPool.PingContext(ctx); err != nil {
		return nil, err
	}

	return &DB{db: connectionPool}, nil
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Close() error {
	return db.db.Close()
}
