package database

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

// TaskRepository の具体的な実装
type TaskRepository struct{}

func (repo *TaskRepository) FindByUser(db *gorm.DB, uid int) (tasks []entity.Task, err error) {
	tasks = []entity.Task{}
	err = db.Where("userid = ?", uid).Find(&tasks).Error
	return
}

func (repo *TaskRepository) FindByID(db *gorm.DB, id int) (task entity.Task, err error) {
	task = entity.Task{}
	err = db.First(&task, id).Error
	return
}

func (repo *TaskRepository) Create(db *gorm.DB, t entity.Task) (task entity.Task, err error) {
	err = db.Create(&t).Error
	return t, err
}

func (repo *TaskRepository) Update(db *gorm.DB, t entity.Task) (task entity.Task, err error) {
	err = db.Save(&t).Error
	return t, err
}

func (repo *TaskRepository) Delete(db *gorm.DB, id int) (taskid int, err error) {
	err = db.Delete(&entity.Task{}, id).Error
	return id, err
}
