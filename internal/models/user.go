package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-"`
	Tasks        []Task `json:"tasks" gorm:"foreignKey:AssigneeID"`
}
