package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t Task) All(db *gorm.DB) ([]Task, error) {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
