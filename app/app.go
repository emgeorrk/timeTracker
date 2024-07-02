package app

import (
	"github.com/jinzhu/gorm"
	"sync"
	"timeTracker/database"
)

type App struct {
	DB        *gorm.DB
	WaitGroup *sync.WaitGroup
}

func NewApp() *App {
	db := database.InitDB()
	
	return &App{DB: db, WaitGroup: &sync.WaitGroup{}}
}
