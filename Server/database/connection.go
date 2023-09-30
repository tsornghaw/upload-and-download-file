package database

import (
	"upload-and-download-file/models"

	"golang.org/x/crypto/bcrypt"
)

var GDB *GormDatabase

const SecretKey = "secret"

func Connect() *GormDatabase {
	// Alvin: Data from config.go
	config := models.DefaultConfig

	GDB, err := NewGormDatabase(
		config.Postgresql.UserName,
		config.Postgresql.Password,
		config.Postgresql.Host,
		config.Postgresql.Port,
		config.Postgresql.DatabaseName,
		config.Postgresql.DatabaseType,
	)

	if err != nil {
		panic("could not connect to the database")
	}

	// CreateTable will create tables if they don't exist
	if err := GDB.CreateTable(&models.User{}); err != nil {
		panic("could not create user table")
	}

	// CreateTable will create tables if they don't exist
	if err := GDB.CreateTable(&models.StroeData{}); err != nil {
		panic("could not create storedata table")
	}

	// Check if adminstrator already existed.
	if err := GDB.GetCorresponding(&models.User{}, "email = ?", config.Admin.Email); err != nil {

		// Create adminstrator
		password, err := bcrypt.GenerateFromPassword(config.Admin.Password, 14)

		if err != nil {
			panic(err)
		}

		adminUser := models.User{
			Name:     config.Admin.Name,
			Email:    config.Admin.Email,
			Password: password,
			Admin:    config.Admin.Admin,
		}

		if err := GDB.Create(&adminUser); err != nil {
			panic("could not create adminstrator")
		}
	}

	return GDB
}
