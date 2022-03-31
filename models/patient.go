package models

import (
	"gorm.io/gorm"
)

type Patient struct {

	//created_at, updated_at, deleted_at
	gorm.Model

	FirstName string        `json:"firstname"`
	LastName  string        `json:"lastname"`
	Email     string        `gorm:"unique_index" json:"email"`
	Password  string        `json:"password,omitempty"`
	Role      string        `json:"role.omitempty"`
	History   []Appointment `gorm:"foreignKey:patient_id" json:"history,omitempty"`
}
