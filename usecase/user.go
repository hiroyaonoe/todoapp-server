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

// UserInteractor は複数のエンティティを操作する際に活用できる
type UserInteractor struct {
	User repository.UserRepository
}

func NewUserInteractor(user repository.UserRepository) *UserInteractor {
	return &UserInteractor{User: user}
}

func (interactor *UserInteractor) Get(id string) (user *entity.User, err error) {
	// User の取得
	user, err = interactor.User.FindByID(id)
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

	// 新規Userを作成
	err = interactor.User.Create(user)

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

	// Userデータを更新
	err = interactor.User.Update(user)

	return
}

func (interactor *UserInteractor) Delete(id string) (err error) {
	// Userデータを削除
	err = interactor.User.Delete(id)
	return
}
