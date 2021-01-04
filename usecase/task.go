package usecase

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
)

type TaskInteractor struct {
	DB   repository.DBRepository
	Task repository.TaskRepository
}

func (interactor *TaskInteractor) Create(task *entity.Task) (err error) {
	// databaseのnot null制約があるので不要？
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if task.Title.IsNull() || task.UserID.IsNull() || task.Date.IsNull() {
		return errs.ErrInvalidTask
	}
	// 不正なユーザーリクエストの判別(TaskIDがnilでない場合)
	if !task.ID.IsNull() {
		return errs.ErrInvalidTask
	}

	// UUIDを付与
	task.NewID()

	db := interactor.DB.Connect()
	// 新規Taskを作成
	err = interactor.Task.Create(db, task)

	return
}

func (interactor *TaskInteractor) GetByID(tid, uid string) (task *entity.Task, err error) {
	db := interactor.DB.Connect()
	// Task の取得
	task, err = interactor.Task.FindByID(db, tid, uid)
	return
}

func (interactor *TaskInteractor) Update(task *entity.Task) (err error) {
	// databaseのnot null制約があるので不要？
	// // 不正なユーザーリクエストの判別(全フィールドがnilの場合)
	// if task.Title.IsNull() && task.Date.IsNull() {
	// 	return errs.ErrInvalidTask
	// }
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if task.Title.IsNull() || task.Date.IsNull() {
		return errs.ErrInvalidTask
	}
	// 不正なユーザーリクエストの判別(UserID or TaskIDがnilの場合)
	if task.UserID.IsNull() || task.ID.IsNull() {
		return errs.ErrInvalidTask
	}

	db := interactor.DB.Connect()
	// Taskデータを更新
	err = interactor.Task.Update(db, task)

	return
}
