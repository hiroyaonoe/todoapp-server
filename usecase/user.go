/*
Package usecase is Application Business Rules.
具体的な処理を行う
domainにのみ依存
*/
package usecase

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
)

type UserInteractor struct {
	DB   repository.DBRepository
	User repository.UserRepository
}

func (interactor *UserInteractor) Get(id int) (jsonUser *entity.UserForJSON, err error) {
	db := interactor.DB.Connect()
	// User の取得
	user, err := interactor.User.FindByID(db, id)
	if err != nil {
		return nil, err
	}
	jsonUser = user.ToUserForJSON()
	return jsonUser, nil
}

func (interactor *UserInteractor) Create(user *entity.User) (jsonUser *entity.UserForJSON, err error) {
	// databaseのnot null制約があるので不要？
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if user.Name.IsNull() || user.Password.IsNull() || user.Email.IsNull() {
		return nil, entity.ErrInvalidUser
	}
	// 不正なユーザーリクエストの判別(UserIDがnilでない場合)
	if user.ID != 0 {
		return nil, entity.ErrInvalidUser
	}

	db := interactor.DB.Connect()
	// 新規Userを作成
	err = interactor.User.Create(db, user)
	if err != nil {
		return nil, err
	}
	jsonUser = user.ToUserForJSON()
	return jsonUser, nil
}

func (interactor *UserInteractor) Update(user *entity.User) (jsonUser *entity.UserForJSON, err error) {
	// 不正なユーザーリクエストの判別(全フィールドがnilの場合)
	if user.Name.IsNull() && user.Password.IsNull() && user.Email.IsNull() {
		return nil, entity.ErrInvalidUser
	}
	// 不正なユーザーリクエストの判別(UserIDがnilの場合)
	if user.ID == 0 {
		return nil, entity.ErrInvalidUser
	}

	db := interactor.DB.Connect()
	// Userデータを更新
	err = interactor.User.Update(db, user)
	if err != nil {
		return nil, err
	}
	jsonUser = user.ToUserForJSON()
	return jsonUser, nil
}
