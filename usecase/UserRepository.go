package usecase

import (
    "github.com/jinzhu/gorm"

    "github.com/hiroyaonoe/todoapp-server/domain"
)

type UserRepository interface {
    FindByID(db *gorm.DB, id int) (user domain.Users, err error)
}