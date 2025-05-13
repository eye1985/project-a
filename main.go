package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"project-a/internal/database"
	"project-a/internal/server"
	"project-a/internal/util"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	envs, err := util.GetEnvs(
		[]string{
			"POSTGRES_URL",
			"PGX5_URL", // Used in makefile
			"WS_URL",
			"HASH_KEY",
			"BLOCK_KEY",
			"ORIGIN",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	wsUrl := envs["WS_URL"]
	hashKey := envs["HASH_KEY"]
	blockKey := envs["BLOCK_KEY"]
	origin := envs["ORIGIN"]
	pgUrl := envs["POSTGRES_URL"]

	pool, dbErr := database.Pool(pgUrl)
	if dbErr != nil {
		fmt.Printf("Error connecting to pool: %v", dbErr)
		os.Exit(1)
	}
	defer pool.Close()
	database.Migrate(pgUrl)

	log.Printf("Connected to PostgreSQL")
	log.Printf("Starting server at port: %v \n", server.PORT)
	log.Fatal(
		server.Serve(
			pool, &server.ServeArgs{
				HashKey:  hashKey,
				BlockKey: blockKey,
				WsUrl:    wsUrl,
				Origin:   origin,
			},
		),
	)
}
