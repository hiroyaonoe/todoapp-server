/*
Package controllers is Interface Adapters.
routerから要求された処理をusecaseにつなぐ

*/
package controllers

import (
	"errors"
	"net/http"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
	"github.com/hiroyaonoe/todoapp-server/usecase"
)

type UserController struct {
	Interactor *usecase.UserInteractor
}

func NewUserController(user repository.UserRepository) *UserController {
	return &UserController{Interactor: usecase.NewUserInteractor(user)}
}

// Get is the Handler for GET /user
func (controller *UserController) Get(c Context) {
	id, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	user, err := controller.Interactor.Get(id)

	if err != nil {
		if errors.Is(err, entity.ErrRecordNotFound) {
			errorToJSON(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Create is the Handler for POST /user
func (controller *UserController) Create(c Context) {
	user, err := getUserFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, ErrBadRequest)
		return
	}

	err = controller.Interactor.Create(user)

	if err != nil {
		var sqlerr *entity.ErrMySQL
		if errors.As(err, &sqlerr) {
			if sqlerr.Number == 0x426 {
				errorToJSON(c, http.StatusBadRequest, ErrDuplicatedEmail)
				return
			}
		}
		if errors.Is(err, usecase.ErrInvalidUser) {
			errorToJSON(c, http.StatusBadRequest, ErrBadRequest)
			return
		}
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update is the Handler for PUT /user
func (controller *UserController) Update(c Context) {
	id, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}
	user, err := getUserFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, ErrBadRequest)
		return
	}
	user.SetID(id)

	err = controller.Interactor.Update(user)

	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUser) {
			errorToJSON(c, http.StatusBadRequest, ErrBadRequest)
			return
		}
		if errors.Is(err, entity.ErrRecordNotFound) {
			errorToJSON(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		var sqlerr *entity.ErrMySQL
		if errors.As(err, &sqlerr) {
			if sqlerr.Number == 0x426 {
				errorToJSON(c, http.StatusBadRequest, ErrDuplicatedEmail)
				return
			}
		}
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update is the Handler for DELETE /user
func (controller *UserController) Delete(c Context) {
	id, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	err = controller.Interactor.Delete(id)

	if err != nil {
		// if errors.Is(err, usecase.ErrInvalidUser) {
		// 	errorToJSON(c, http.StatusBadRequest, ErrBadRequest)
		// 	return
		// }
		if errors.Is(err, entity.ErrRecordNotFound) {
			errorToJSON(c, http.StatusNotFound, ErrUserNotFound)
			return
		}
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func getUserFromBody(c Context) (user *entity.User, err error) {
	err = c.ShouldBindJSON(&user)
	return
}
