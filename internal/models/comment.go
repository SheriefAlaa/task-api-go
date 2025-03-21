package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Comment string `json:"comment"`
	TaskID  uint   `json:"task_id" gorm:"index"`
	Task    Task   `json:"task" gorm:"foreignKey:TaskID;"`
	UserID  uint   `json:"user_id" gorm:"index"`
	User    User   `json:"user" gorm:"foreignKey:UserID;"`
}
