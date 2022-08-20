package handler

import (
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/model"

	"github.com/gin-gonic/gin"
)

type TaskHander struct{}

type Task model.Task

func (t TaskHander) Index(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": "task not found"})
		return
	}

	c.IndentedJSON(200, tasks)
}

func (t TaskHander) Show(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	id := c.Params.ByName("id")

	var task Task
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, task)
}

func (t TaskHander) Create(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := db.Create(&task).Error; err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(201, task)
}

func (t TaskHander) Update(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	id := c.Params.ByName("id")

	var task Task
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := db.Save(&task).Error; err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(200, task)
}

func (t TaskHander) Delete(c *gin.Context) {
	sqlDB, db := client.ProvidePostgreSqlClient()
	defer sqlDB.Close()

	id := c.Params.ByName("id")

	var task Task
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	if err := db.Delete(&task).Error; err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(204, gin.H{"message": "deleted"})
}
