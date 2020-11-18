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
	err = tx.First(&beforeuser, u.ID).Error
	if err != nil {
		tx.Rollback()
		return
	}
	FillInNilFields(beforeuser, u)
	err = tx.Save(u).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (repo *UserRepository) Delete(db *gorm.DB, id int) (err error) {
	tx := db.Begin()

	// idに該当するユーザーがいない場合を弾く
	user := &entity.User{}
	err = tx.First(user, id).Error
	if err != nil {
		tx.Rollback()
		return
	}

	// err = tx.Delete(&entity.User{}, id).Error
	err = tx.Where("id = ?", id).Delete(&entity.User{}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func FillInNilFields(before entity.User, after *entity.User) {
	if after.Name.IsNull() {
		after.Name = before.Name
	}
	if after.Password.IsNull() {
		after.Password = before.Password
	}
	if after.Email.IsNull() {
		after.Email = before.Email
	}
	return
}
