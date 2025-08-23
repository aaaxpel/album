package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func connect() *pgx.Conn {
	_ = godotenv.Load()

	if url, ok := os.LookupEnv("DB"); ok {
		conn, err := pgx.Connect(context.Background(), url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		defer conn.Close(context.Background())

		err = conn.Ping(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Pong")
		}

		return conn

	}
	return nil
}
