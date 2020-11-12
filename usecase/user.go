/*
usecase is Application Business Rules.
具体的な処理を行う
*/
package usecase

import (
	"net/http"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
)

type UserInteractor struct {
	DB   repository.DBRepository
	User repository.UserRepository
}

func (interactor *UserInteractor) Get(id int) (user entity.User, resultStatus *ResultStatus) {
	db := interactor.DB.Connect()
	// User の取得
	user, err := interactor.User.FindByID(db, id)
	if err == entity.ErrRecordNotFound {
		return entity.User{}, NewResultStatus(http.StatusNotFound, entity.ErrUserNotFound)
	}
	if err != nil {
		return entity.User{}, NewResultStatus(http.StatusInternalServerError, err)
	}
	(&user).HidePassword()
	return user, NewResultStatus(http.StatusOK, nil)
}

func (interactor *UserInteractor) Create(user *entity.User) (resultStatus *ResultStatus) {
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if (user.Name == "") || (user.Password == "") || (user.Email == "") {
		return NewResultStatus(http.StatusBadRequest, entity.ErrInvalidUser)
	}

	db := interactor.DB.Connect()
	// 新規Userを作成
	err := interactor.User.Create(db, user)
	if err != nil {
		return NewResultStatus(http.StatusInternalServerError, err)
	}
	user.HidePassword()
	return NewResultStatus(http.StatusOK, nil)
}

func (interactor *UserInteractor) Update(user *entity.User) (resultStatus *ResultStatus) {
	// 不正なユーザーリクエストの判別(全フィールドがnilの場合)
	if (user.Name == "") && (user.Password == "") && (user.Email == "") {
		return NewResultStatus(http.StatusBadRequest, entity.ErrInvalidUser)
	}

	db := interactor.DB.Connect()
	// Userデータを更新
	err := interactor.User.Update(db, user)
	if err != nil {
		return NewResultStatus(http.StatusInternalServerError, err)
	}
	user.HidePassword()
	return NewResultStatus(http.StatusOK, nil)
}
