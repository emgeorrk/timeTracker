package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	PassportSeries string `json:"passport_series"`
	PassportNumber string `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	Tasks          []Task `json:"tasks"`
}

type Task struct {
	gorm.Model
	UserID           uint          `json:"user_id"`
	Name             string        `json:"name"`
	IsActive         bool          `json:"is_active"`
	Periods          []Period      `json:"periods"`
	OverallTimeSpent time.Duration `json:"overall_time"`
}

type Period struct {
	gorm.Model
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
