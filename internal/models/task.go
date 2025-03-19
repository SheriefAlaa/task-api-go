package models

import "gorm.io/gorm"

// Task represents a task entity in our system
type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssigneeID  uint   `json:"assignee_id"`
	Assignee    User   `json:"assignee" gorm:"foreignKey:AssigneeID"`
}
