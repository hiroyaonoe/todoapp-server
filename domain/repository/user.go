//go:generate mockgen -source=$GOFILE -destination=../mock_repository/mock_$GOFILE -package=mock_repository

/*
Package repository is Enterprise Business Rules.
データベースへの処理がどうあるべきかインターフェースの形で記述
永続化を責務とする
どこにも依存しない
*/
package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
)

// UserRepository is interface of User
type UserRepository interface {
	FindByID(db *gorm.DB, uuid int) (user *entity.User, err error)
	Create(db *gorm.DB, u *entity.User) (err error)
	Update(db *gorm.DB, u *entity.User) (err error)
	Delete(db *gorm.DB, uuid int) (err error)
}
