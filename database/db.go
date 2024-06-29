package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"timeTracker/config"
	"timeTracker/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("DB_HOST"), config.GetEnv("DB_USER"), config.GetEnv("DB_PASSWORD"), config.GetEnv("DB_NAME"), config.GetEnv("DB_PORT"))
	
	DB, err = gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	
	DB.AutoMigrate(&models.User{}, &models.Task{})
}
