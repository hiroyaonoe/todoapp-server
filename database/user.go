/*
Package database is Interface Adapters.
SQLへのクエリはここで行う
*/
package database

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

// UserRepository の具体的な実装
type UserRepository struct{}

func (repo *UserRepository) FindByID(db *gorm.DB, id int) (user *entity.User, err error) {
	user = &entity.User{}
	err = db.First(user, id).Error
	return
}

func (repo *UserRepository) Create(db *gorm.DB, u *entity.User) (err error) {
	tx := db.Begin()
	err = tx.Create(u).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (repo *UserRepository) Update(db *gorm.DB, u *entity.User) (err error) {
	tx := db.Begin()
	beforeuser := entity.User{}
	err = db.First(&beforeuser, u.ID).Error
	if err != nil {
		tx.Rollback()
		return
	}
	FillInNullFields(beforeuser, u)
	err = db.Save(u).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (repo *UserRepository) Delete(db *gorm.DB, id int) (uid int, err error) {
	tx := db.Begin()
	err = tx.Delete(&entity.User{}, id).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return id, err
}

func FillInNullFields(before entity.User, after *entity.User) {
	if after.Name == "" {
		after.Name = before.Name
	}
	if after.Password == "" {
		after.Password = before.Password
	}
	if after.Email == "" {
		after.Email = before.Email
	}
	return
}
