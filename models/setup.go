package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details

	dsn := "rb:cxunicorn123@@tcp(cx-go-tas-rb.mysql.database.azure.com)/HOSPITAL?charset=utf8mb4&parseTime=True"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database!")
	}

	//Migrating my schema to keep it up to date.
	database.AutoMigrate(&Doctor{}, &Appointment{})
	database.AutoMigrate(&Patient{}, &Appointment{})
	database.AutoMigrate(&Admin{})
	DB = database
}
