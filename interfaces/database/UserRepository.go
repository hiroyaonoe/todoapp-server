package database

import (
    "errors"

    "github.com/hiroyaonoe/todoapp-server/domain"
)

type UserRepository struct {}

func (repo *UserRepository) FindByID(db *gorm.DB, id int) (user domain.User, err error) {
    user = domain.User{}
    db.First(&user, id)
    if user.ID <= 0 {
        return domain.User{}, errors.New("user is not found")
    }
    return user, nil
}