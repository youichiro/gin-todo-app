package main

import (
	"example/web-service-gin/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	router.GET("/tasks", handler.TaskHander{}.Index)
	router.Run("0.0.0.0:8080")
}
