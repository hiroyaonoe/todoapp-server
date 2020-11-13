/*
Package controllers is Interface Adapters.
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

// Get is the Handler for GET /user
func (controller *UserController) Get(c Context) {
	id, err := GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewH(err.Error(), nil))
		return
	}

	jsonUser, res := controller.Interactor.Get(id)
	if res.Error != nil {
		c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
		return
	}
	c.JSON(res.StatusCode, NewH("success", jsonUser))
}

// Create is the Handler for POST /user
func (controller *UserController) Create(c Context) {
	user, err := GetUserFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewH(err.Error(), nil))
		return
	}

	jsonUser, res := controller.Interactor.Create(&user)
	if res.Error != nil {
		c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
		return
	}
	c.JSON(res.StatusCode, NewH("success", jsonUser))
}

// Update is the Handler for PUT /user
func (controller *UserController) Update(c Context) {
	user, err := GetUserFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewH(err.Error(), nil))
		return
	}

	id, err := GetUserIDFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewH(err.Error(), nil))
		return
	}
	user.ID = id

	jsonUser, res := controller.Interactor.Update(&user)
	if res.Error != nil {
		c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
		return
	}
	c.JSON(res.StatusCode, NewH("success", jsonUser))
}

func GetUserIDFromCookie(c Context) (id int, err error) {
	cookie, err := c.Cookie("id")
	id, err = strconv.Atoi(cookie)
	return
}

func GetUserFromBody(c Context) (user entity.User, err error) {
	err = c.ShouldBindJSON(&user)
	return
}
