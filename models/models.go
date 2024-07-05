package models

import (
	"time"
)

type User struct {
	ID             uint   `gorm:"primary_key" json:"id"`
	PassportSeries string `json:"passport_series"`
	PassportNumber string `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	Tasks          []Task `gorm:"foreignKey:UserID" json:"tasks"`
}

type Task struct {
	ID               uint          `gorm:"primary_key;autoIncrement:false" json:"id"`
	UserID           uint          `gorm:"primary_key;autoIncrement:false" json:"user_id"`
	Name             string        `json:"name"`
	IsActive         bool          `json:"is_active"`
	Periods          []Period      `gorm:"foreignKey:ID;References:ID,TaskID" json:"periods"`
	OverallTimeSpent time.Duration `json:"overall_time"`
}

type Period struct {
	TaskID    uint      `gorm:"primary_key;autoIncrement:false" json:"task_id"`
	UserID    uint      `gorm:"primary_key;autoIncrement:false" json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
