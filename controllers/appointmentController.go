package controllers

import (
	"net/http"
	"strconv"
	"techunicorn/models"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Booking struct {
	PatientID uint          `json:"patientID"`
	DoctorID  uint          `json:"doctorID"`
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Duration  time.Duration `json:"duration"`
}

// POST function to book an appointment
// /doctors/:id/book
func CreateAppointment(c *gin.Context) {
	var input Booking
	var doc models.Doctor
	var count int64
	var totTime float32
	var overlap bool
	var patient models.Patient

	doctorID, er := strconv.Atoi(c.Param("id"))
	doc.ID = uint(doctorID)

	email := c.GetString("email")

	if err := models.DB.Where("email = ?", email).First(&patient).Error; err != nil {

		//Only patients can book an appointment otherwise, throws an error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "WRONG USER!"})
		return
	}

	if er != nil {
		c.JSON(http.StatusBadRequest, nil)
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attempt := input.StartTime
	attempt2 := input.EndTime
	now := time.Now()

	//Checking if appointment date is not in the past
	if attempt.Before(now) {

		//If it is, then throws an error message
		c.JSON(http.StatusBadRequest, gin.H{"Error": "You cannot book an appointment in the past :)"})
		return
	}

	// Minimum Duration: 15 mins, Max Duration: 2 Hours
	//Duration of the appointment
	if input.EndTime.Sub(input.StartTime).Minutes() < 15 || input.EndTime.Sub(input.StartTime).Hours() > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Duration of the appointment must be between 15 minutes to 2 hours"})
		return
	}

	//Reference: stackoverflow.com/questions/44873825/how-to-get-timestamp-of-utc-time-with-golang

	t := attempt.Format("2017-07-02")

	models.DB.Raw("SELECT COUNT(*) FROM hospital.appointments WHERE DATE(start_time) = ? AND d_id=?",
		string(t), doc.ID).Scan(&count)

	// a doctor can have a maximum of 12 appointments and not more than that
	if count >= 12 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "The doctor already has 12 patients for this day! Please choose another day."})
		return
	}

	//Reference: stackoverflow.com/questions/4102480/mysql-how-to-sum-a-timediff-on-a-group

	models.DB.Raw("SELECT SUM((TIMEDIFF(end_time, start_time)/10000)) FROM hospital.appointments WHERE DATE(start_time) = ? AND d_id=?",
		string(t), doc.ID).Scan(&totTime)

	// a doctor can have a maximum total appointments time of 8 hours in a day
	if totTime >= 8 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "A doctor can have a maximum total appointments time of 8 hours in a day! Book another day!"})
		return
	}

	//Checking that given booking is after 9 AM and before 5 PM of that date
	year, month, day := attempt.Date()

	//Reference: geeksforgeeks.org/time-time-utc-function-in-golang-with-examples/

	// Defining t for UTC method
	// t := time.Date(2020, 11, 14, 11, 30, 32, 0, time.UTC)

	time9AM := time.Date(year, month, day, 9, 0, 0, 0, time.UTC)
	time5PM := time.Date(year, month, day, 17, 0, 0, 0, time.UTC)

	if attempt.Before(time9AM) || attempt2.After(time5PM) {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Appointments are available only during 9 AM to 5 PM"})
		return
	}

	//Reference: stackoverflow.com/questions/54618633/parsing-time-as-2006-01-02t150405z0700-cannot-parse-as-2006

	//panic: parsing time """" as ""2006-01-02T15:04:05Z07:00"

	//Checking for duplicate time
	parsingtime1 := attempt.Format("2006-01-02T15:04:05Z07:00")
	parsingtime2 := attempt2.Format("2006-01-02T15:04:05Z07:00")

	//Arbitrary SQL commands which aren't parsed any further by the query builder.
	//They therefore can create a vector for attack via SQL injection.
	models.DB.Raw("SELECT EXISTS (SELECT * FROM hospital.appointments WHERE doctor_id=? AND start_time BETWEEN ? AND ? OR end_time BETWEEN ? and ?)",
		doc.ID, string(parsingtime1), string(parsingtime2), string(parsingtime1), string(parsingtime2)).Scan(&overlap)
	if overlap {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your appointment time overlaps with another appointment!"})
		return
	}

	//Checking if doctor is available or not
	if err := models.DB.First(&doc).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Doctor not found!"})
		return
	}

	//When all conditions are good, an appointment is created
	appotmt := models.Appointment{
		PatientID: patient.ID,
		DoctorID:  doc.ID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Duration:  time.Duration(input.EndTime.Sub(input.StartTime).Minutes())}
	models.DB.Create(&appotmt)

	c.JSON(http.StatusOK, gin.H{"Appointment is booked!!": appotmt})
	//Go through the key instead------
	//Return it as appointment itself, easier to parse------
}

// DELETE function to cancel an appointment
// /appointments/:id
func DeleteAppointment(d *gin.Context) {
	role := d.GetString("role")

	if role != "admin" && role != "doctor" {
		d.JSON(http.StatusBadRequest, gin.H{"Error": "Doctors only can cancel appointments"})
		return
	}
	var doctor models.Doctor

	//strconv.Atoi is used to convert string type into int type.
	aID, er := strconv.Atoi(d.Param("id"))
	if er != nil {
		d.JSON(http.StatusBadRequest, nil)
	}

	var appt models.Appointment
	if err := models.DB.Where("id = ?", aID).First(&appt).Error; err != nil {
		d.JSON(http.StatusBadRequest, gin.H{"Error": "Record not found!"})
		return
	}

	if role == "doctor" {
		email := d.GetString("email")
		models.DB.Where("email = ?", email).First(&doctor)
		if doctor.ID != appt.DoctorID {
			d.JSON(http.StatusBadRequest, gin.H{"Error": "Incorrect doctor! Booked Doctor only can cancel this appointment"})
			return
		}
	}
	models.DB.Delete(&appt)
	d.JSON(http.StatusOK, gin.H{"Appointments is cancelled": true})
}

//GET function to view all the appointments in the system
// /appointments/:id
func ViewAppointment(v *gin.Context) {
	role := v.GetString("role")

	var appointment models.Appointment
	var doctor models.Doctor
	var patient models.Patient

	if err := models.DB.Where("id = ?", v.Param("id")).First(&appointment).Error; err != nil {
		v.JSON(http.StatusBadRequest, gin.H{"Error": "Record not found!"})
		return
	}

	if role == "doctor" {
		models.DB.Where("email = ?", v.GetString("email")).First(&patient)
		if patient.ID == appointment.PatientID {

			//Raw SQL database queries
			models.DB.Raw("SELECT first_name, last_name, email FROM hospital.patients WHERE id=?", appointment.PatientID).Find(&patient)
			models.DB.Raw("SELECT first_name, last_name, email FROM hospital.doctors WHERE id=?", appointment.DoctorID).Find(&doctor)
			v.JSON(http.StatusOK, gin.H{"appointment": appointment, "patient": patient, "doctor": doctor})
		} else {
			v.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect doctor information!"})
		}

	} else if role == "patient" {
		models.DB.Where("email = ?", v.GetString("email")).First(&doctor)
		if doctor.ID == appointment.DoctorID {

			//Raw SQL database queries
			models.DB.Raw("SELECT first_name, last_name, email FROM hospital.patients WHERE id=?", appointment.PatientID).Find(&patient)
			models.DB.Raw("SELECT first_name, last_name, email FROM hospital.doctors WHERE id=?", appointment.DoctorID).Find(&doctor)
			v.JSON(http.StatusOK, gin.H{"appointment": appointment, "patient": patient, "doctor": doctor})
		} else {
			v.JSON(http.StatusBadRequest, gin.H{"Error": "Incorrect patient information!"})
		}

		//Only clinic admin can see the appointments

	} else {
		//Raw SQL database queries
		models.DB.Raw("SELECT first_name, last_name, email FROM hospital.patients WHERE id=?", appointment.PatientID).Find(&patient)
		models.DB.Raw("SELECT first_name, last_name, email FROM hospital.doctors WHERE id=?", appointment.DoctorID).Find(&doctor)
		v.JSON(http.StatusOK, gin.H{"appointment": appointment, "patient": patient, "doctor": doctor})
	}
}
