package main

import (
	"example/web-service-gin/internal/client"
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

	db := client.PostgresClientProvider{}
	db.Connect()
	defer db.Close()

	router := gin.Default()
	router.GET("/tasks", handler.TaskHander{}.Index)
	router.GET("/tasks/:id", handler.TaskHander{}.Show)
	router.POST("/tasks", handler.TaskHander{}.Create)
	router.PUT("/tasks/:id", handler.TaskHander{}.Update)
	router.DELETE("/tasks/:id", handler.TaskHander{}.Delete)

	router.Run("0.0.0.0:8080")
}
