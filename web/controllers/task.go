package controllers

import (
	"net/http"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
	"github.com/hiroyaonoe/todoapp-server/usecase"
)

type TaskController struct {
	Interactor usecase.TaskInteractor
}

func NewTaskController(db repository.DBRepository, task repository.TaskRepository) *TaskController {
	return &TaskController{
		Interactor: usecase.TaskInteractor{
			DB:   db,
			Task: task,
		},
	}
}

// Create is the Handler for POST /task
func (controller *TaskController) Create(c Context) {
	task, err := getTaskFromBody(c)
	if err != nil {
		errorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
		return
	}

	userid, err := getUserIDFromCookie(c)
	task.UserID.Set(userid)

	err = controller.Interactor.Create(task)

	if err != nil {
		if err == entity.ErrInvalidTask {
			errorToJSON(c, http.StatusBadRequest, entity.ErrBadRequest)
			return
		}
		panic(err.Error())
		// errorToJSON(c, http.StatusInternalServerError, entity.ErrInternalServerError)
		// return
	}
	c.JSON(http.StatusOK, task)
}

func getTaskFromBody(c Context) (task *entity.Task, err error) {
	err = c.ShouldBindJSON(&task)
	return
}
