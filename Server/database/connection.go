package database

import (
	"encoding/json"
	"fmt"
	"upload-and-download-file/models"
)

var GDB *GormDatabase

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

	PrettyPrint(config.Postgresql.UserName)
	PrettyPrint(config.Postgresql.Password)
	PrettyPrint(config.Postgresql.Host)
	PrettyPrint(config.Postgresql.Port)
	PrettyPrint(config.Postgresql.DatabaseName)
	PrettyPrint(config.Postgresql.DatabaseType)

	if err != nil {
		panic("could not connect to the database")
	}

	if GDB == nil {
		panic("fail to connect database")
	} else if GDB != nil {
		fmt.Printf("successfully connect to postgres ...\n")
	}

	//DB = connection
	// CreateTable will create tables if they don't exist
	GDB.CreateTable(&models.User{})
	if err != nil {
		panic("could not create table")
	}

	return GDB
}

// print the contents of the obj
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
