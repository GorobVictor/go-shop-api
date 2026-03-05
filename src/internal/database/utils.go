package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection() (context.Context, *pgxpool.Pool, error) {
	connStr := os.Getenv("POSTGRES_CONN_STR")

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, connStr)

	return ctx, conn, err
}
