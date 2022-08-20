package handler

import (
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHander struct{}

type Task model.Task

func (t TaskHander) Index(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, tasks)
}
