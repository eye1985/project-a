package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"project-a/server"
)

func main() {
	pgUrl, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		pgUrl = "postgres://admin:root@postgres:5432/project_a"
	}

	conn, err := pgx.Connect(context.Background(), pgUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	log.Printf("Connected to PostgreSQL")
	log.Printf("Starting server at port: %v \n", server.PORT)
	log.Fatal(server.Serve())
}
