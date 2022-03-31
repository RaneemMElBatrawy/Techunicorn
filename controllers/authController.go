package controllers

import (
	"log"
	"net/http"

	"techunicorn/middleware"
	"techunicorn/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//Registering Structure /register
type SignUp struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

//JSON Web Token structure 1) Header 2) Payload 3) Signature
//Payload contains the claims.
//Claims are statements about an entity (typically, the user) and additional data.
//There are three types of claims: registered, public, and private claims.

//Reference:
//https://gist.github.com/hamzawix/56f5120fa23295f39ccd4a081e694afc

// LoginPayload login body
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role" `
}

// LoginResponse token response
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// Signup creates a user in db - REGISTERING
func Signup(s *gin.Context) {
	var input SignUp
	err := s.ShouldBindJSON(&input)
	if err != nil {
		s.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing field!"})
		s.Abort()
		return
	}

	//Reference: stackoverflow.com/questions/69432893/crypto-bcrypt-hashedpassword-is-not-the-hash-of-the-given-password

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 8)

	//DOCTOR-----------------
	if input.Role == "doctor" {
		doctor := models.Doctor{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  string(hashed),
			Role:      input.Role}

		//Check if the doctor has already created an email in the clinic's system
		var doc models.Doctor
		if err := models.DB.Where("email = ?", doctor.Email).First(&doc).Error; err == nil {
			s.JSON(http.StatusBadRequest, gin.H{"Error": "As a doctor, your email already exists!"})
			return
		}
		models.DB.Create(&doctor)
		s.JSON(http.StatusOK, doctor)
	}

	//PATIENT------------------
	if input.Role == "patient" {
		patient := models.Patient{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  string(hashed),
			Role:      input.Role}

		var pat models.Patient
		if err := models.DB.Where("email = ?", patient.Email).First(&pat).Error; err == nil {
			s.JSON(http.StatusBadRequest, gin.H{"Error": "As a Patient, your email already exists!"})
			return
		}

		models.DB.Create(&patient)
		s.JSON(http.StatusOK, patient)

	}

	//ADMIN------------------
	if input.Role == "admin" {
		admin := models.Admin{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  string(hashed),
			Role:      input.Role}

		var ad models.Admin
		if err := models.DB.Where("email = ?", admin.Email).First(&ad).Error; err == nil {
			s.JSON(http.StatusBadRequest, gin.H{"Error": "As an admin, your email already exists!"})
			return
		}

		models.DB.Create(&admin)
		s.JSON(http.StatusOK, admin)
	}
}

// Login logs users in - LOGING IN
func Login(l *gin.Context) {
	var payload LoginPayload
	var user models.Doctor
	var doctor models.Doctor
	var patient models.Patient
	var admin models.Admin

	err := l.ShouldBindJSON(&payload)
	if err != nil {
		l.JSON(http.StatusBadRequest, gin.H{
			"Error": "One of the fields is missing. Invalid JSON!!",
		})
		l.Abort()
		return
	}
	if payload.Role == "doctor" {
		result := models.DB.Where("email = ?", payload.Email).Find(&(doctor))
		if result.Error == gorm.ErrRecordNotFound {
			l.JSON(400, gin.H{
				"Error": "Invalid user credentials!!",
			})
			l.Abort()
			return
		}
		copier.CopyWithOption(&user, &doctor, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	}

	if payload.Role == "patient" {
		result := models.DB.Where("email = ?", payload.Email).Find(&patient)
		if result.Error == gorm.ErrRecordNotFound {
			l.JSON(401, gin.H{
				"Error": "Invalid user credentials!!",
			})
			l.Abort()
			return
		}
		copier.CopyWithOption(&user, &patient, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	}
	if payload.Role == "admin" {
		result := models.DB.Where("email = ?", payload.Email).Find(&admin)
		if result.Error == gorm.ErrRecordNotFound {
			l.JSON(401, gin.H{
				"Error": "Invalid user credentials!!",
			})
			l.Abort()
			return
		}
		copier.CopyWithOption(&user, &admin, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	}

	//Reference: gowebexamples.com/password-hashing/

	e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if e != nil {
		l.JSON(401, gin.H{
			"Error": "Wrong password!!",
		})
		l.Abort()
		return
	}

	//Reference: ichi.pro/de/hands-on-mit-jwt-in-golang-154586293449920

	jwtWrapper := middleware.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email, user.Role)
	if err != nil {
		log.Println(err)
		l.JSON(500, gin.H{
			"messsage": "error signing token",
		})
		l.Abort()
		return
	}

	tokenResponse := LoginResponse{
		Token:   signedToken,
		Message: "Login Successful! Woho!",
	}

	l.JSON(http.StatusOK, tokenResponse)
}
