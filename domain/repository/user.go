package usecase

import (
    "github.com/jinzhu/gorm"

    "github.com/hiroyaonoe/todoapp-server/domain"
)

type UserRepository interface {
    FindByID(db *gorm.DB, id int) (user domain.User, err error)
    Create(db *gorm.DB, u domain.User) (user domain.User, err error)
    Update(db *gorm.DB, u domain.User) (user domain.User, err error)
    Delete(db *gorm.DB, id int) (uid int, err error)
}