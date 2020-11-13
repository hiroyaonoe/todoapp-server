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
	err = db.Create(u).Error
	return
}

func (repo *UserRepository) Update(db *gorm.DB, u *entity.User) (err error) {
	beforeuser := entity.User{}
	err = db.First(&beforeuser, u.ID).Error
	if err != nil {
		return err
	}
	FillInNullFields(beforeuser, u)
	err = db.Save(u).Error
	return
}

func (repo *UserRepository) Delete(db *gorm.DB, id int) (uid int, err error) {
	err = db.Delete(&entity.User{}, id).Error
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
