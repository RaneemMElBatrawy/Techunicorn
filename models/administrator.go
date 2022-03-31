package models

import (
	"gorm.io/gorm"
)

type Admin struct {

	//created_at, updated_at, deleted_at
	gorm.Model

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"type"`
}
