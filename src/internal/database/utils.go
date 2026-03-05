package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection() (context.Context, *pgxpool.Pool, error) {
	connStr := os.Getenv("POSTGRES_CONN_STR")

	if connStr == "" {
		log.Fatal("connection string is empty")
	}

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, connStr)

	return ctx, conn, err
}
