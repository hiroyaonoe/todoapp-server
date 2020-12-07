/*
Package controllers is Interface Adapters.
routerから要求された処理をusecaseにつなぐ

*/
package controllers

import (
	"net/http"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
	"github.com/hiroyaonoe/todoapp-server/usecase"
)

type UserController struct {
	Interactor usecase.UserInteractor
}

func NewUserController(db repository.DBRepository, user repository.UserRepository) *UserController {
	return &UserController{
		Interactor: usecase.UserInteractor{
			DB:   db,
			User: user,
		},
	}
}

// Get is the Handler for GET /user
func (controller *UserController) Get(c Context) {
	id, err := getUserIDFromCookie(c)
	if err != nil {
		ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		return
	}

	jsonUser, err := controller.Interactor.Get(id)

	if err != nil {
		if err == entity.ErrRecordNotFound {
			ErrorToJSON(c, http.StatusNotFound, entity.ErrUserNotFound)
			return
		}
		ErrorToJSON(c, http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, jsonUser)
}

// Create is the Handler for POST /user
func (controller *UserController) Create(c Context) {
	user, err := getUserFromBody(c)
	if err != nil {
		ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		return
	}

	jsonUser, err := controller.Interactor.Create(&user)

	if err != nil {
		if err == entity.ErrInvalidUser {
			ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
			return
		}
		ErrorToJSON(c, http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, jsonUser)
}

// Update is the Handler for PUT /user
func (controller *UserController) Update(c Context) {
	user, err := getUserFromBody(c)
	if err != nil {
		ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		return
	}
	id, err := getUserIDFromCookie(c)
	if err != nil {
		ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		return
	}
	user.SetID(id)

	jsonUser, err := controller.Interactor.Update(&user)

	if err != nil {
		if err == entity.ErrInvalidUser {
			ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
			return
		}
		if err == entity.ErrRecordNotFound {
			ErrorToJSON(c, http.StatusNotFound, entity.ErrUserNotFound)
			return
		}
		ErrorToJSON(c, http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, jsonUser)
}

func (controller *UserController) Delete(c Context) {
	id, err := getUserIDFromCookie(c)
	if err != nil {
		ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		return
	}

	err = controller.Interactor.Delete(id)

	if err != nil {
		// if err == entity.ErrInvalidUser {
		// 	ErrorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		// 	return
		// }
		if err == entity.ErrRecordNotFound {
			ErrorToJSON(c, http.StatusNotFound, entity.ErrUserNotFound)
			return
		}
		ErrorToJSON(c, http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func getUserFromBody(c Context) (user entity.User, err error) {
	err = c.ShouldBindJSON(&user)
	return
}
