package main

import (
	"github.com/youichiro/go-todo-app/internal/client"
	"github.com/youichiro/go-todo-app/internal/router"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/volatiletech/sqlboiler/boil"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("GO_ENV")

	boil.DebugMode = true // なぜか効かない
	client.Connect(env)
	defer client.DB.Close()

	r := router.SetupRouter()
	r.Run("0.0.0.0:8080")
}
