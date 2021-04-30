package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `validate:"required,email" gorm:"not null;unique"`
	Name     string `validate:"required,min=2,max=100" gorm:"not null"`
	Password string `validate:"required,min=8,max=100" gorm:"not null"`
	Roles    string `validate:"required,min=2,max=500" gorm:"not null;default:'user'"`
	Active   bool   `validate:"required" gorm:"not null;default:true"`
}

type Pagination struct {
	Sort string `json:"sort"`
	Page int    `json:"page"`
	Size int    `json:"size"`
}
