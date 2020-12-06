package usecase

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/repository"
)

type TaskInteractor struct {
	DB   repository.DBRepository
	Task repository.TaskRepository
}

func (interactor *TaskInteractor) Create(task *entity.Task) (jsonTask *entity.TaskForJSON, err error) {
	// databaseのnot null制約があるので不要？
	// 不正なユーザーリクエストの判別(フィールドのうち少なくともひとつがnilの場合)
	if task.Title.IsNull() || task.UserID.IsNull() || task.Date.IsNull() {
		return nil, entity.ErrInvalidUser
	}
	// 不正なユーザーリクエストの判別(TaskIDがnilでない場合)
	if !task.ID.IsNull() {
		return nil, entity.ErrInvalidUser
	}

	// UUIDを付与
	task.NewID()

	db := interactor.DB.Connect()
	// 新規Userを作成
	err = interactor.Task.Create(db, task)
	if err != nil {
		return nil, err
	}
	jsonTask = task.ToTaskForJSON()
	return jsonTask, nil
}
