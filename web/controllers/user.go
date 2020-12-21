/*
Package controllers is Interface Adapters.
routerから要求された処理をusecaseにつなぐ

*/
package controllers

import (
	"net/http"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
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
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	user, err := controller.Interactor.Get(id)

	if err != nil {
		if err == errs.ErrRecordNotFound {
			errorToJSON(c, http.StatusNotFound, errs.ErrUserNotFound)
			return
		}
		panic(err.Error())
		// errorToJSON(c, http.StatusInternalServerError, errs.ErrInternalServerError)
		// return
	}
	c.JSON(http.StatusOK, user)
}

// Create is the Handler for POST /user
func (controller *UserController) Create(c Context) {
	user, err := getUserFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	err = controller.Interactor.Create(user)

	if err != nil {
		if e, ok := err.(*errs.ErrMySQL); ok {
			if e.Number == 0x426 {
				errorToJSON(c, http.StatusBadRequest, errs.ErrDuplicatedEmail)
				return
			}
		}
		if err == errs.ErrInvalidUser {
			errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}
		panic(err.Error())
		// errorToJSON(c, http.StatusInternalServerError, errs.ErrInternalServerError)
		// return
	}
	c.JSON(http.StatusOK, user)
}

// Update is the Handler for PUT /user
func (controller *UserController) Update(c Context) {
	user, err := getUserFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}
	id, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}
	user.SetID(id)

	err = controller.Interactor.Update(user)

	if err != nil {
		if err == errs.ErrInvalidUser {
			errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}
		if err == errs.ErrRecordNotFound {
			errorToJSON(c, http.StatusNotFound, errs.ErrUserNotFound)
			return
		}
		panic(err.Error())
		// errorToJSON(c, http.StatusInternalServerError, errs.ErrInternalServerError)
		// return
	}
	c.JSON(http.StatusOK, user)
}

func (controller *UserController) Delete(c Context) {
	id, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	err = controller.Interactor.Delete(id)

	if err != nil {
		// if err == errs.ErrInvalidUser {
		// 	errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		// 	return
		// }
		if err == errs.ErrRecordNotFound {
			errorToJSON(c, http.StatusNotFound, errs.ErrUserNotFound)
			return
		}
		panic(err.Error())
		// errorToJSON(c, http.StatusInternalServerError, errs.ErrInternalServerError)
		// return
	}
	c.JSON(http.StatusOK, nil)
}

func getUserFromBody(c Context) (user *entity.User, err error) {
	err = c.ShouldBindJSON(&user)
	return
}
