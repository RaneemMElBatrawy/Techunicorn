package main

import (
	"techunicorn/controllers"
	"techunicorn/middleware"
	"techunicorn/models"

	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()
	models.ConnectDatabase()

	// Routes
	public := route.Group("/public")
	register := route.Group("/")

	register.Use(middleware.Auth()) //JWT bearer token

	//1) /register - allow all users to register
	public.POST("/register", controllers.Signup)

	//2) /login - allow users to login
	public.POST("/login", controllers.Login)

	//3) /doctors - view lists of doctors
	register.GET("/doctors", controllers.ViewAllDoctors)

	//4) /doctors/:id - view doctors information
	register.GET("/doctors/:id", controllers.ViewRequestedDoctor)

	//5) /doctors/:id/slots - view doctor available slots
	register.POST("/doctors/:id/slots", controllers.ViewDoctorSlots)

	//6) /doctors/:id/book - book an appointment with a doctor
	register.POST("/doctors/:id/book", controllers.CreateAppointment)

	//7) /appointments/:id - cancel appointment
	register.DELETE("/appointments/:id", controllers.DeleteAppointment)

	//8) /doctors//availability/all - view availability of all doctors
	register.POST("/doctors/availability/all", controllers.ViewDoctorAvailability)

	//9) /appointments/:id - view appointments details
	register.GET("/appointments/:id", controllers.ViewAppointment)

	// /patients - View all the patients in the system
	register.GET("/patients", controllers.ViewAllPatients)

	// /patients/:id - View all the requested patients with their id in the system
	register.GET("/patients/:id", controllers.ViewRequestedPatient)

	//10) /patients/:id/history - view patient appointment history
	register.GET("/patients/:id/history", controllers.ViewPatientHistory)

	//11) /doctors/most/appointments - view doctors with most appointments in a given day
	register.POST("/doctors/most/appointments", controllers.ViewDoctorsMostAppointments)

	//12) /doctors/most/hours - view doctors who have 6+ hours total appointment in a day
	register.POST("/doctors/most/hours", controllers.ViewDoctorsMostHours)

	//Running the server
	route.Run()

}
