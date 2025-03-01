package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"project-a/database"
	"project-a/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	pool, dbErr := database.Pool()
	if dbErr != nil {
		fmt.Printf("Error connecting to pool: %v", dbErr)
		os.Exit(1)
	}
	defer pool.Close()
	database.Migrate()

	log.Printf("Connected to PostgreSQL")
	log.Printf("Starting server at port: %v \n", server.PORT)
	log.Fatal(server.Serve(pool))
}
