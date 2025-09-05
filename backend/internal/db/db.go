package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() *pgxpool.Pool {
	_ = godotenv.Load()

	if url, ok := os.LookupEnv("DB"); ok {
		pool, err := pgxpool.New(context.Background(), url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		return pool

	}
	return nil
}
