/*
Package usecase is Application Business Rules.
具体的な処理を行う
domainにのみ依存
*/
package usecase

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
)

type UserInteractor struct {
	DB   repository.DBRepository
	User repository.UserRepository
}

func (interactor *UserInteractor) Get(id string) (user *entity.User, err error) {
	db := interactor.DB.Connect()
	// User の取得
	user, err = interactor.User.FindByID(db, id)

	return
}

func (interactor *UserInteractor) Create(user *entity.User) (err error) {
	// databaseのnot null制約があるので不要？
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if user.Name.IsNull() || user.Password.IsNull() || user.Email.IsNull() {
		return errs.ErrInvalidUser
	}
	// 不正なユーザーリクエストの判別(UserIDがnilでない場合)
	if !user.ID.IsNull() {
		return errs.ErrInvalidUser
	}

	// UUIDを付与
	user.NewID()

	// Passwordをhash化
	EncryptPassword(user)

	db := interactor.DB.Connect()
	// 新規Userを作成
	err = interactor.User.Create(db, user)

	return
}

func (interactor *UserInteractor) Update(user *entity.User) (err error) {
	// 不正なユーザーリクエストの判別(全フィールドがnilの場合)
	if user.Name.IsNull() && user.Password.IsNull() && user.Email.IsNull() {
		return errs.ErrInvalidUser
	}
	// 不正なユーザーリクエストの判別(UserIDがnilの場合)
	if user.ID.IsNull() {
		return errs.ErrInvalidUser
	}

	// Passwordをhash化
	EncryptPassword(user)

	db := interactor.DB.Connect()
	// Userデータを更新
	err = interactor.User.Update(db, user)

	return
}

func (interactor *UserInteractor) Delete(id string) (err error) {
	db := interactor.DB.Connect()
	// Userデータを削除
	err = interactor.User.Delete(db, id)

	return
}
