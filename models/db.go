package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func NewRepositoryFromEnvironments() *Repository {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	if username == "" || password == "" || dbName == "" || dbHost == "" {
		return nil
	}
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	//fmt.Println(dbUri)

	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Connected to database")
	db := conn
	db.AutoMigrate(&Offer{})

	return &Repository{db: db}
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *Repository) SetDB(gdb *gorm.DB) {
	r.db = gdb
}
