// @title Time Tracker API
// @version 1.0
// @description Time tracking application.
// @host localhost:8084
// @BasePath /api/v1/

package main

import (
	"log"
	"timeTracker/app"
	"timeTracker/config"
	"timeTracker/external"
	"timeTracker/routes"
	"timeTracker/swagger"
)

func main() {
	config.LoadConfig()
	
	myApp := app.NewApp()
	
	defer func(myApp *app.App) {
		err := myApp.DB.Close()
		if err != nil {
			log.Fatalln("Failed to close database connection:", err)
		}
	}(myApp)
	
	myApp.WaitGroup.Add(1)
	go external.RunExternalApiEmulation(myApp)
	myApp.WaitGroup.Wait()
	
	go swagger.RunSwagger()
	
	r := routes.InitRouter(myApp)
	
	log.Println("Server started on http://localhost:8084")
	if err := r.Run(":8084"); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
