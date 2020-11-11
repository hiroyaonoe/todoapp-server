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

func NewUserController(db database.DBRepository) *UserController {
	return &UserController{
		Interactor: usecase.UserInteractor{
			DB:   &db,
			User: &database.UserRepository{},
		},
	}
}

func (controller *UserController) Get(c Context) error {
	cookie, _ := c.Cookie("id")
	id, _ := strconv.Atoi(cookie)

	user, res, err := controller.Interactor.Get(id)
	if err != nil {
		c.JSON(res.StatusCode, res)
		return err
	}
	c.JSON(res.StatusCode, user)
	return nil
}
