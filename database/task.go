package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

// TaskRepository の具体的な実装
type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{db: db.Connect()}
}

func (repo *TaskRepository) Create(t *entity.Task) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr)
		}
	}()

	tx := repo.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	t.NewID()

	err = tx.Create(t).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (repo *TaskRepository) FindByID(tid, uid string) (task *entity.Task, err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr) //TODO:testなし
		}
		return
	}()

	task = &entity.Task{}
	err = repo.db.Where("id = ?", tid).Where("user_id = ?", uid).First(task).Error
	return
}

func (repo *TaskRepository) Update(t *entity.Task) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr) //TODO:testなし
		}
	}()

	tx := repo.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	task := &entity.Task{}
	// idに該当するユーザーがいない場合を弾く
	err = tx.Where("id = ?", t.ID).Where("user_id = ?", t.UserID).First(task).Error
	if err != nil {
		return
	}

	err = tx.Save(t).Error
	if err != nil {
		return //TODO:testなし
	}
	return
}

func (repo *TaskRepository) Delete(tid, uid string) (err error) {
	defer func() {
		if nerr, ok := err.(*mysql.MySQLError); ok {
			err = (*entity.ErrMySQL)(nerr) //TODO:testなし
		}
	}()

	tx := repo.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	task := &entity.Task{}
	// idに該当するユーザーがいない場合を弾く
	err = tx.Where("id = ?", tid).Where("user_id = ?", uid).First(task).Error
	if err != nil {
		return
	}
	err = repo.db.Where("id = ?", tid).Where("user_id = ?", uid).Delete(&entity.Task{}).Error
	if err != nil {
		return //TODO:testなし
	}
	return
}
