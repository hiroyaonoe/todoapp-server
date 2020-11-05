package database

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

// UserRepository の具体的な実装
type UserRepository struct{}

func (repo *UserRepository) FindByID(db *gorm.DB, id int) (user entity.User, err error) {
	user = entity.User{}
	err = db.First(&user, id).Error
	return
}

func (repo *UserRepository) Create(db *gorm.DB, u entity.User) (user entity.User, err error) {
	err = db.Create(&u).Error
	return u, err
}

func (repo *UserRepository) Update(db *gorm.DB, u entity.User) (user entity.User, err error) {
	err = db.Save(&u).Error
	return u, err
}

func (repo *UserRepository) Delete(db *gorm.DB, id int) (uid int, err error) {
	err = db.Delete(&entity.User{}, id).Error
	return id, err
}
