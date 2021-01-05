package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/jinzhu/gorm"
)

// TaskRepository の具体的な実装
type TaskRepository struct{}

func (repo *TaskRepository) Create(db *gorm.DB, t *entity.Task) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr)
		}
	}()

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Create(t).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (repo *TaskRepository) FindByID(db *gorm.DB, tid, uid string) (task *entity.Task, err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr) //TODO:testなし
		}
		return
	}()

	task = &entity.Task{}
	err = db.Where("id = ?", tid).Where("user_id = ?", uid).First(task).Error
	return
}

func (repo *TaskRepository) Update(db *gorm.DB, t *entity.Task) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*errs.ErrMySQL)(nerr) //TODO:testなし
		}
	}()

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	beforetask := entity.Task{}
	err = tx.Where("id = ?", t.ID).Where("user_id = ?", t.UserID).First(&beforetask).Error
	if err != nil {
		return
	}
	// FillInTaskNilFields(beforetask, t)
	err = tx.Save(t).Error
	if err != nil {
		return //TODO:testなし
	}
	return
}

// func (repo *TaskRepository) FindByUser(db *gorm.DB, uid int) (tasks []entity.Task, err error) {
// 	tasks = []entity.Task{}
// 	err = db.Where("userid = ?", uid).Find(&tasks).Error
// 	return
// }

// func (repo *TaskRepository) Delete(db *gorm.DB, id int) (taskid int, err error) {
// 	err = db.Delete(&entity.Task{}, id).Error
// 	return id, err
// }

// func FillInTaskNilFields(before entity.Task, after *entity.Task) {
// 	if after.Title.IsNull() {
// 		after.Title = before.Title
// 	}
// 	// TODO: Contentを更新しないのか，空に更新したいのかが判別不能
// 	if after.Content.IsNull() {
// 		after.Content = before.Content
// 	}
// 	// TODO: IsCompを更新しないのか，falseに更新したいのかが判別不能
// 	if after.Date.IsNull() {
// 		after.Date = before.Date
// 	}
// 	return
// }
