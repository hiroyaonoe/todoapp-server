package usecase

import (
    "github.com/hiroyaonoe/todoapp-server/domain"
)

type UserInteractor struct {
    DB DBRepository
    User UserRepository
}

func (interactor *UserInteractor) Get(id int) (user domain.UserForGet, resultStatus *ResultStatus) {
    db := interactor.DB.Connect()
    // User の取得
    foundUser, err := interactor.User.FindByID(db, id)
    if err != nil {
        return domain.UserForGet{}, NewResultStatus(404, err)
    }
    user = foundUser.BuildForGet()
    return user, NewResultStatus(200, nil)
}