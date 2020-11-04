package usecase

import (
    "github.com/jinzhu/gorm"

    "github.com/hiroyaonoe/todoapp-server/domain"
)

type TaskRepository interface {
    FindByUser(db *gorm.DB, uid int) (tasks []domain.Task, err error)
    FindByID(db *gorm.DB, id int) (task domain.Task, err error)
    Create(db *gorm.DB, t domain.Task) (task domain.Task, err error)
    Update(db *gorm.DB, t domain.Task) (task domain.Task, err error)
    Delete(db *gorm.DB, id int) (taskid int, err error)
}