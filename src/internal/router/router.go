package router

import (
	"database/sql"

	"github.com/youichiro/go-todo-app/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	taskHandler := handler.TaskHander{DB: db}

	r.GET("/", func(c *gin.Context) { c.IndentedJSON(200, gin.H{"message": "hello world!"}) })
	r.GET("/tasks", taskHandler.Index)
	r.GET("/tasks/:id", taskHandler.Show)
	r.POST("/tasks", taskHandler.Create)
	r.PUT("/tasks/:id", taskHandler.Update)
	r.DELETE("/tasks/:id", taskHandler.Delete)

	return r
}
