package main

import (
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := client.PostgresClientProvider{}
	db.Connect()
	defer db.Close()

	r := router.SetupRouter()
	r.Run("0.0.0.0:8080")
}
