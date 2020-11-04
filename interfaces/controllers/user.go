/*
controllers is Interface Adapters.
routerから要求された処理をusecaseにつなぐ
*/
package controllers

import (
    "strconv"

    "github.com/hiroyaonoe/todoapp-server/interfaces/database"
    "github.com/hiroyaonoe/todoapp-server/usecase"
)

type UserController struct {
    Interactor usecase.UserInteractor
}

func NewUserController(db database.DB) *UserController {
    return &UserController{
        Interactor: usecase.UserInteractor{
            DB: &database.DBRepository{ DB: db },
            User: &database.UserRepository{},
        },
    }
}

func (controller *UserController) Get(c Context) {

    id, _ := strconv.Atoi(c.Param("id"))

    user, res := controller.Interactor.Get(id)
    if res.Error != nil {
        c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
        return
    }
    c.JSON(res.StatusCode, NewH("success", user))
}