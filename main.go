// @title Time Tracker API
// @version 1.0
// @description Time tracking application.
// @host localhost:8084
// @BasePath /

package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"timeTracker/config"
	"timeTracker/database"
	"timeTracker/routes"
)

func main() {
	config.LoadConfig()
	
	database.InitDB()
	defer func(DB *gorm.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalln("Failed to close database connection:", err)
		}
	}(database.DB)
	
	r := routes.InitRouter()
	
	err := r.Run("localhost:8084")
	if err != nil {
		log.Fatalln("Failed to start server:", err)
	}
	
	log.Println("Server started on http://localhost:8084")
}
