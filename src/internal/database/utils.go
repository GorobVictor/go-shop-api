package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection() (*pgxpool.Pool, error) {
	connStr := os.Getenv("POSTGRES_CONN_STR")

	if connStr == "" {
		log.Fatal("connection string is empty")
	}

	conn, err := pgxpool.New(context.Background(), connStr)

	return conn, err
}
