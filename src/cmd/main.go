package main

import (
	"example/web-service-gin/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/tasks", handler.TaskIndex)
	router.Run("0.0.0.0:8080")
}
