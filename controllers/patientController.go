package controllers

import (
	"fmt"
	"net/http"

	"techunicorn/models"

	"github.com/gin-gonic/gin"
)

//Reference: https://stackoverflow.com/questions/38501646/initialize-nested-struct-definition-in-golang-if-it-have-same-objects

type CreatePatient struct {
	FirstName string               `json:"firstname"`
	LastName  string               `json:"lastname"`
	Email     string               `json:"email" `
	Password  string               `json:"password"`
	History   []models.Appointment `json:"history"`
}
type DoctorToPatient struct {
	FirstName string
	LastName  string
	Email     string
}
type DoctorToPatientHistory struct {
	Appointment models.Appointment
	Doctor      DoctorToPatient
}

//GET function to view all the patients in the system
// /patients
func ViewAllPatients(v *gin.Context) {
	role := v.GetString("role")
	fmt.Println(role)
	if role == "admin" {
		var patients []models.Patient
		models.DB.Find(&patients)

		v.JSON(http.StatusOK, gin.H{"patients": patients})
	} else {
		v.JSON(http.StatusBadRequest, gin.H{"error": "NOT A USER! Out of scope!"})
	}
}

//GET function to view all the requested patients with their id in the system
// /patients/:id
func ViewRequestedPatient(r *gin.Context) {

	var doc models.Patient
	role := r.GetString("role")
	if role == "admin" {
		if err := models.DB.Where("id = ?", r.Param("id")).First(&doc).Error; err != nil {

			//If the given the id is not in the record, throw an error message

			r.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		r.JSON(http.StatusOK, gin.H{"patient": doc})
	} else {
		//If not an admin, then throw a message, only admins can view this list
		r.JSON(http.StatusBadRequest, gin.H{"error": "NOT A USER! Out of scope!"})
	}
}

//GET function to view the history of the patients
// /patients/:id/history
func ViewPatientHistory(h *gin.Context) {
	var patient models.Patient
	var pat models.Patient
	var doc models.Doctor
	history := []DoctorToPatientHistory{}

	email := h.GetString("email")
	role := h.GetString("role")

	if role == "patient" {
		models.DB.Where("email=?", email).First(&pat)
		if fmt.Sprint(pat.ID) != h.Param("id") {

			//Patient can only view their own history and cannot view history of other patients for privacy
			h.JSON(http.StatusBadRequest, gin.H{"error": "Cannot to be accessed!"})
			return
		}
	} else {
		if err := models.DB.Where("id = ?", h.Param("id")).First(&patient).Error; err != nil {

			//If the patient's id is not in the record, throw a message
			h.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}

		models.DB.Preload("History").Find(&patient)

		for i := 0; i < len(patient.History); i++ {
			take := patient.History[i].DoctorID

			models.DB.Raw("SELECT * FROM hospital.doctors WHERE id=?", take).First(&doc)
			doctor := DoctorToPatient{
				doc.FirstName,
				doc.LastName,
				doc.Email}

			data := DoctorToPatientHistory{patient.History[i], doctor}

			history = append(history, data)
		}
		h.JSON(http.StatusOK, gin.H{"history": history})
	}
}
