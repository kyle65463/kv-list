package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DbPool *pgxpool.Pool

func CreateDbConnection() {
	dbUrl := os.Getenv("DATABASE_URL")

	var err error
	DbPool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't connect to the database: %v\n", err)
		os.Exit(1)
	}
}
