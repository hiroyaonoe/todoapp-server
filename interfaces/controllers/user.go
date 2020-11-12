/*
controllers is Interface Adapters.
routerから要求された処理をusecaseにつなぐ
*/
package controllers

import (
	"net/http"
	"strconv"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
	"github.com/hiroyaonoe/todoapp-server/interfaces/database"
	"github.com/hiroyaonoe/todoapp-server/usecase"
)

type UserController struct {
	Interactor usecase.UserInteractor
}

func NewUserController(db repository.DBRepository) *UserController {
	return &UserController{
		Interactor: usecase.UserInteractor{
			DB:   db,
			User: &database.UserRepository{},
		},
	}
}

func (controller *UserController) Get(c Context) {
	cookie, err := c.Cookie("id")
	id, err := strconv.Atoi(cookie)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewH(err.Error(), nil))
		return
	}

	user, res := controller.Interactor.Get(id)
	if res.Error != nil {
		c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
		return
	}
	c.JSON(res.StatusCode, NewH("success", user))
}

func (controller *UserController) Create(c Context) {
	user := entity.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewH(err.Error(), nil))
		return
	}

	fields := []string{user.Name, user.Password, user.Email}
	for _, field := range fields {
		if field == "" {
			c.JSON(http.StatusBadRequest, NewH("invalid user", nil))
			return
		}
	}

	res := controller.Interactor.Create(&user)
	if res.Error != nil {
		c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
		return
	}
	c.JSON(res.StatusCode, NewH("success", user))
}
