/*
Package database is Interface Adapters.
SQLへのクエリはここで行う
*/
package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

// UserRepository の具体的な実装
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db.Connect()}
}

func (repo *UserRepository) FindByID(id string) (user *entity.User, err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr) //TODO:testなし
		}
		return
	}()
	user = &entity.User{}
	err = repo.db.Where("id = ?", id).First(user).Error
	return
}

func (repo *UserRepository) Create(u *entity.User) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr)
		}
		return
	}()

	tx := repo.db.Begin()
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

func (repo *UserRepository) Update(u *entity.User) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr)
		}
		return
	}()

	tx := repo.db.Begin()
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

func (repo *UserRepository) Delete(id string) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr) //TODO:testなし
		}
		return
	}()
	tx := repo.db.Begin()
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
