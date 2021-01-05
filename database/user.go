/*
Package database is Interface Adapters.
SQLへのクエリはここで行う
*/
package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/jinzhu/gorm"
)

// UserRepository の具体的な実装
type UserRepository struct{}

func (repo *UserRepository) FindByID(db *gorm.DB, id string) (user *entity.User, err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr) //TODO:testなし
		}
		return
	}()
	user = &entity.User{}
	err = db.Where("id = ?", id).First(user).Error
	return
}

func (repo *UserRepository) Create(db *gorm.DB, u *entity.User) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr)
		}
	}()

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Create(u).Error
	if err != nil {
		return
	}
	return
}

func (repo *UserRepository) Update(db *gorm.DB, u *entity.User) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr)
		}
		return
	}()

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	beforeuser := entity.User{}
	err = tx.Where("id = ?", u.ID).First(&beforeuser).Error
	if err != nil {
		return
	}
	FillInUserNilFields(beforeuser, u)
	err = tx.Save(u).Error
	if err != nil {
		return
	}
	return
}

func (repo *UserRepository) Delete(db *gorm.DB, id string) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr) //TODO:testなし
		}
		return
	}()
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// idに該当するユーザーがいない場合を弾く
	user := &entity.User{}
	err = tx.Where("id = ?", id).First(user).Error
	if err != nil {
		return
	}

	// err = tx.Delete(&entity.User{}, id).Error
	err = tx.Where("id = ?", id).Delete(&entity.User{}).Error
	if err != nil {
		return //TODO:testなし
	}
	return
}

func FillInUserNilFields(before entity.User, after *entity.User) {
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
