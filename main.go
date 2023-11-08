package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/R894/go-blog/api"
	"github.com/R894/go-blog/models"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db, err := models.NewPostgresDatabase()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	server := api.NewServer(port, logger, db)
	log.Fatal(server.Start())
}
