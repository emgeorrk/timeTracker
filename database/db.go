package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"timeTracker/config"
	"timeTracker/models"
)

func InitDB() *gorm.DB {
	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("DB_HOST"), config.GetEnv("DB_USER"), config.GetEnv("DB_PASSWORD"), config.GetEnv("DB_NAME"), config.GetEnv("DB_PORT"))
	
	DB, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	log.Println("Database connection established")
	
	if err := DB.AutoMigrate(&models.User{}, &models.Task{}, &models.Period{}).Error; err != nil {
		log.Fatalln("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
	
	return DB
}
