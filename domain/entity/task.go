package entity

import (
	"time"
)

// Task は内部で処理する際のTask情報である
type Task struct {
	ID          int    `gorm:"primary_key"`
	Title       string `gorm:"not null"`
	Content     string
	UserID      int `gorm:"not null;index"`
	IsCompleted bool
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
