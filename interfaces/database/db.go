package database

import (
    "github.com/jinzhu/gorm"
)

type DBRepository struct {
}

func (db *DBRepository) Begin() *gorm.DB {
    return db.Begin()
}

func (db *DBRepository) Connect() *gorm.DB {
    return db.Connect()
}