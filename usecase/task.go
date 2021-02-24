package usecase

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
)

// TaskInteractor は複数のエンティティを操作する際に活用できる
type TaskInteractor struct {
	Task repository.TaskRepository
}

func NewTaskInteractor(task repository.TaskRepository) *TaskInteractor {
	return &TaskInteractor{Task: task}
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

	// 新規Taskを作成
	err = interactor.Task.Create(task)

	return
}

func (interactor *TaskInteractor) GetByID(tid, uid string) (task *entity.Task, err error) {
	// Task の取得
	task, err = interactor.Task.FindByID(tid, uid)
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

	// Taskデータを更新
	err = interactor.Task.Update(task)

	return
}

func (interactor *TaskInteractor) Delete(tid, uid string) (err error) {
	// Taskの削除
	err = interactor.Task.Delete(tid, uid)
	return
}
