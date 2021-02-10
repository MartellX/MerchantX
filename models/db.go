package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func init() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	if username == "" || password == "" || dbName == "" || dbHost == "" {
		return
	}
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	//fmt.Println(dbUri)

	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Connected to database")
	db = conn
	db.AutoMigrate(&Offer{})
}

func GetDB() *gorm.DB {
	return db
}

func SetDB(gdb *gorm.DB)  {
	db = gdb
}
