/*
Package usecase is Application Business Rules.
具体的な処理を行う
domainにのみ依存
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

func (interactor *UserInteractor) Get(id int) (jsonUser entity.UserForJSON, resultStatus *ResultStatus) {
	db := interactor.DB.Connect()
	// User の取得
	user, err := interactor.User.FindByID(db, id)
	if err == entity.ErrRecordNotFound {
		return entity.UserForJSON{}, NewResultStatus(http.StatusNotFound, entity.ErrUserNotFound)
	}
	if err != nil {
		return entity.UserForJSON{}, NewResultStatus(http.StatusInternalServerError, err)
	}
	jsonUser = user.ToUserForJSON()
	return jsonUser, NewResultStatus(http.StatusOK, nil)
}

func (interactor *UserInteractor) Create(user *entity.User) (jsonUser entity.UserForJSON, resultStatus *ResultStatus) {
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if (user.Name == "") || (user.Password == "") || (user.Email == "") {
		return entity.UserForJSON{}, NewResultStatus(http.StatusBadRequest, entity.ErrInvalidUser)
	}

	db := interactor.DB.Connect()
	// 新規Userを作成
	err := interactor.User.Create(db, user)
	if err != nil {
		return entity.UserForJSON{}, NewResultStatus(http.StatusInternalServerError, err)
	}
	jsonUser = user.ToUserForJSON()
	return jsonUser, NewResultStatus(http.StatusOK, nil)
}

func (interactor *UserInteractor) Update(user *entity.User) (jsonUser entity.UserForJSON, resultStatus *ResultStatus) {
	// 不正なユーザーリクエストの判別(全フィールドがnilの場合)
	if (user.Name == "") && (user.Password == "") && (user.Email == "") {
		return entity.UserForJSON{}, NewResultStatus(http.StatusBadRequest, entity.ErrInvalidUser)
	}

	db := interactor.DB.Connect()
	// Userデータを更新
	err := interactor.User.Update(db, user)
	if err == entity.ErrRecordNotFound {
		return entity.UserForJSON{}, NewResultStatus(http.StatusNotFound, entity.ErrUserNotFound)
	}
	if err != nil {
		return entity.UserForJSON{}, NewResultStatus(http.StatusInternalServerError, err)
	}
	jsonUser = user.ToUserForJSON()
	return jsonUser, NewResultStatus(http.StatusOK, nil)
}
