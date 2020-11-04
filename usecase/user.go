/*
usecase is Application Business Rules.
*/
package usecase

import (
    "github.com/hiroyaonoe/todoapp-server/domain/entity"
    "github.com/hiroyaonoe/todoapp-server/domain/repository"
)

type UserInteractor struct {
    DB repository.DBRepository
    User repository.UserRepository
}

func (interactor *UserInteractor) Get(id int) (user entity.UserForGet, resultStatus *ResultStatus) {
    db := interactor.DB.Connect()
    // User の取得
    foundUser, err := interactor.User.FindByID(db, id)
    if err != nil {
        return entity.UserForGet{}, NewResultStatus(404, err)
    }
    user = foundUser.BuildForGet()
    return user, NewResultStatus(200, nil)
}