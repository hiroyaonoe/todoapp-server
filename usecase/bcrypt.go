package usecase

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(u *entity.User) error {
	password := []byte(u.Password.GetString())
	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = entity.NewNullString(string(hashed))

	return nil
}

func ComparePassword(hashed *entity.User, plain *entity.User) error {
	p1 := []byte(hashed.Password.GetString())
	p2 := []byte(plain.Password.GetString())

	return bcrypt.CompareHashAndPassword(p1, p2)
}
