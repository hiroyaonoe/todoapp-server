package database

import (
    "github.com/jinzhu/gorm"
    "github.com/hiroyaonoe/todoapp-server/domain"
)

// UserRepository の具体的な実装
type UserRepository struct {}

func (repo *UserRepository) FindByID(db *gorm.DB, id int) (user domain.User, err error) {
    user = domain.User{}
    err = db.First(&user, id).Error
    return
}

func (repo *UserRepository) Create(db *gorm.DB, u domain.User) (user domain.User, err error) {
    err = db.Create(&u).Error
    return u, err
}

func (repo *UserRepository) Update(db *gorm.DB, u domain.User) (user domain.User, err error) {
    err = db.Save(&u).Error
    return u, err
}

func (repo *UserRepository) Delete(db *gorm.DB, id int) (uid int, err error) {
    err = db.Delete(&domain.User{}, id).Error
    return id, err
}
