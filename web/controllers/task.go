package controllers

import (
	"errors"
	"net/http"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
	"github.com/hiroyaonoe/todoapp-server/usecase"
)

type TaskController struct {
	Interactor *usecase.TaskInteractor
}

func NewTaskController(task repository.TaskRepository) *TaskController {
	return &TaskController{Interactor: usecase.NewTaskInteractor(task)}
}

// Create is the Handler for POST /task
func (controller *TaskController) Create(c Context) {
	uid, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, errs.ErrUnauthorized)
		return
	}
	task, err := getTaskFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	task.UserID.Set(uid)

	err = controller.Interactor.Create(task)

	if err != nil {
		if errors.Is(err, errs.ErrInvalidTask) {
			errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}
		// TODO:user not found
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// GetByID is the Handler for GET /task/:id
func (controller *TaskController) GetByID(c Context) {
	uid, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, errs.ErrUnauthorized)
		return
	}
	tid, err := getTaskIDFromParam(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	task, err := controller.Interactor.GetByID(tid, uid)

	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			errorToJSON(c, http.StatusNotFound, errs.ErrTaskNotFound)
			return
		}
		// TODO:user not found
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// Update is the Handler for PUT /task/:id
func (controller *TaskController) Update(c Context) {
	uid, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, errs.ErrUnauthorized)
		return
	}
	task, err := getTaskFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	tid, err := getTaskIDFromParam(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	task.ID.Set(tid)
	task.UserID.Set(uid)

	err = controller.Interactor.Update(task)

	if err != nil {
		if errors.Is(err, errs.ErrInvalidTask) {
			errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}
		if errors.Is(err, errs.ErrRecordNotFound) {
			errorToJSON(c, http.StatusNotFound, errs.ErrTaskNotFound)
			return
		}
		// TODO:user not found
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// Delete is the Handler for DELETE /task/:id
func (controller *TaskController) Delete(c Context) {
	uid, err := getUserIDFromCookie(c)
	if err != nil {
		errorToJSON(c, http.StatusUnauthorized, errs.ErrUnauthorized)
		return
	}
	tid, err := getTaskIDFromParam(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, errs.ErrBadRequest)
		return
	}

	err = controller.Interactor.Delete(tid, uid)

	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			errorToJSON(c, http.StatusNotFound, errs.ErrTaskNotFound)
			return
		}
		// TODO:user not found
		unexpectedErrorHandling(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func getTaskFromBody(c Context) (task *entity.Task, err error) {
	err = c.ShouldBindJSON(&task)
	return
}
