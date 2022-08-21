package router

import (
	"example/web-service-gin/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) { c.IndentedJSON(200, gin.H{"message": "hello world!"}) })
	r.GET("/tasks", handler.TaskHander{}.Index)
	r.GET("/tasks/:id", handler.TaskHander{}.Show)
	r.POST("/tasks", handler.TaskHander{}.Create)
	r.PUT("/tasks/:id", handler.TaskHander{}.Update)
	r.DELETE("/tasks/:id", handler.TaskHander{}.Delete)

	return r
}
