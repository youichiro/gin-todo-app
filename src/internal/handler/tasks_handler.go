package handler

import (
	"github.com/youichiro/go-todo-app/internal/client"
	"github.com/youichiro/go-todo-app/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TaskHander struct{}

type createParams struct {
	Title string `json:"title" binding:"required"`
}

type updateParams struct {
	Title string `json:"title" binding:"required"`
	Done bool `json:"done"`
}

func (t TaskHander) Index(c *gin.Context) {
	tasks, err := models.Tasks().All(c, client.DB)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(200, tasks)
}

func (t TaskHander) Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	task, err := models.FindTask(c, client.DB, id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(200, task)
}

func (t TaskHander) Create(c *gin.Context) {
	var params createParams
	err := c.BindJSON(&params)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}

	task := &models.Task{
		Title: params.Title,
		Done:  false,
	}

	err = task.Insert(c, client.DB, boil.Infer())
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(201, task)
}

func (t TaskHander) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))

	var params updateParams
	err := c.BindJSON(&params)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": err.Error()})
		return
	}

	task, err := models.FindTask(c, client.DB, id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	task.Title = params.Title
	task.Done = params.Done
	_, err = task.Update(c, client.DB, boil.Infer())
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(200, gin.H{"message": "successfully updated"})
}

func (t TaskHander) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))

	task, err := models.FindTask(c, client.DB, id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	_, err = task.Delete(c, client.DB)
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(204, gin.H{"message": "successfully deleted"})
}
