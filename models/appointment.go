package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {

	//created_at, updated_at, deleted_at
	gorm.Model

	PatientID uint          `json:"patientID"`
	DoctorID  uint          `json:"doctorID"`
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"dateTime"`
	Duration  time.Duration `json:"duration"`
}
