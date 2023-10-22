package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type DB struct {
	db *sql.DB
}

func NewDB(dbname string) (*DB, error) {
	connectionString := fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres dbname=%s sslmode=disable", dbname)

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
