package database

import (
    "github.com/jinzhu/gorm"
    "github.com/hiroyaonoe/todoapp-server/domain"
)

// TaskRepository の具体的な実装
type TaskRepository struct {}

func (repo *TaskRepository) FindByUser(db *gorm.DB, uid int) (tasks []domain.Task, err error) {
    tasks = []domain.Task{}
    err = db.Where("userid = ?", uid).Find(&tasks).Error
    return
}

func (repo *TaskRepository) FindByID(db *gorm.DB, id int) (task domain.Task, err error) {
    task = domain.Task{}
    err = db.First(&task, id).Error
    return
}

func (repo *TaskRepository) Create(db *gorm.DB, t domain.Task) (task domain.Task, err error) {
    err = db.Create(&t).Error
    return t, err
}

func (repo *TaskRepository) Update(db *gorm.DB, t domain.Task) (task domain.Task, err error) {
    err = db.Save(&t).Error
    return t, err
}

func (repo *TaskRepository) Delete(db *gorm.DB, id int) (taskid int, err error) {
    err = db.Delete(&domain.Task{}, id).Error
    return id, err
}
