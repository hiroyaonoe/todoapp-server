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
    DB repository.DBRepository
    User repository.UserRepository
}

func (interactor *UserInteractor) Get(id int) (user entity.User, resultStatus *ResultStatus) {
    db := interactor.DB.Connect()
    // User の取得
    foundUser, err := interactor.User.FindByID(db, id)
    if err != nil {
        return entity.User{}, NewResultStatus(http.StatusNotFound, err)
    }
    user = foundUser.BuildForGet()
    return user, NewResultStatus(http.StatusOK, nil)
}