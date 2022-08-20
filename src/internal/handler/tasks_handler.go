package handler

import (
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TaskIndex(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	tasks, err := model.Task{}.All(db)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, tasks)
}
