package controllers

import (
    "strconv"

    "github.com/hiroyaonoe/todoapp-server/interfaces/database"
    "github.com/hiroyaonoe/todoapp-server/usecase"
)

type TaskController struct {
    Interactor usecase.TaskInteractor
}

func NewTaskController(db database.DB) *TaskController {
    return &TaskController{
        Interactor: usecase.TaskInteractor{
            DB: &database.DBRepository{ DB: db },
            Task: &database.TaskRepository{},
        },
    }
}

func (controller *TaskController) Get(c Context) {

    id, _ := strconv.Atoi(c.Param("id"))

    Task, res := controller.Interactor.Get(id)
    if res.Error != nil {
        c.JSON(res.StatusCode, NewH(res.Error.Error(), nil))
        return
    }
    c.JSON(res.StatusCode, NewH("success", Task))
}