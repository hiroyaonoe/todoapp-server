package domain

import (
	"time"
)

type Task struct{
	ID int `gorm:"primary_key"`
	Title string `sql:"not null"`
	Content string
	UserID int `sql:"not null"`
	IsCompleted bool 
	Date time.Time
	CreatedAt time.Time
    UpdatedAt time.Time
}