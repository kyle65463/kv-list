package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxInterface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Close(context.Context) error
}

var DbPool PgxInterface

func CreateDbConnection() {
	dbUrl := os.Getenv("DATABASE_URL")

	var err error
	DbPool, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't connect to the database: %v\n", err)
		os.Exit(1)
	}
}

func CloseDbConnection() {
	DbPool.Close(context.Background())
}
