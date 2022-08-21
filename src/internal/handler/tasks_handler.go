package handler

import (
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/model"

	"github.com/gin-gonic/gin"
)

type TaskHander struct{}

type Task model.Task

func (t TaskHander) Index(c *gin.Context) {
	var tasks []Task
	if err := client.DB.Find(&tasks).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": "task not found"})
		return
	}

	c.IndentedJSON(200, tasks)
}

func (t TaskHander) Show(c *gin.Context) {
	id := c.Params.ByName("id")
	var task Task

	if err := client.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, task)
}

func (t TaskHander) Create(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := client.DB.Create(&task).Error; err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(201, task)
}

func (t TaskHander) Update(c *gin.Context) {
	id := c.Params.ByName("id")
	var task Task

	if err := client.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := client.DB.Save(&task).Error; err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(200, task)
}

func (t TaskHander) Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	var task Task

	if err := client.DB.Where("id = ?", id).First(&task).Error; err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	if err := client.DB.Delete(&task).Error; err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(204, gin.H{"message": "deleted"})
}
