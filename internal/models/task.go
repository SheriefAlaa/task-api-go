package models

import (
	"database/sql/driver"
	"errors"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusTodo        TaskStatus = "todo"
	StatusInProgress  TaskStatus = "in_progress"
	StatusUnderReview TaskStatus = "under_review"
	StatusDone        TaskStatus = "done"
)

func (s TaskStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *TaskStatus) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("invalid data for TaskStatus")
	}

	switch TaskStatus(str) {
	case StatusTodo, StatusInProgress, StatusUnderReview, StatusDone:
		*s = TaskStatus(str)
		return nil
	default:
		return errors.New("invalid status value")
	}
}

type Task struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status" gorm:"type:string;default:'todo'"`
	AssigneeID  uint       `json:"assignee_id" gorm:"index"`
	Assignee    User       `json:"assignee" gorm:"foreignKey:AssigneeID"`
	Comments    []Comment  `json:"comments" gorm:"foreignKey:TaskID"`
}
