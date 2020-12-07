package database

import (
	"reflect"
	"testing"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

const (
	uuidTA1 = "cb8aa4fc-1964-a965-00ea-8b158c0ffcc7"
	uuidTA2 = "38397cad-8865-081f-3482-2a035f875d5c"
	uuidTB1 = "2265150f-a3c9-d21e-ee75-a8c42a07807e"
	uuidTB2 = "b67a753e-90c0-4435-ede7-fdcda1504ca5"
)

var (
	taskA1 = entity.NewTask(uuidTA1, "taskA1", "I am ContentA1.", uuidUA, "2020-12-08")
	taskA2 = entity.NewTask(uuidTA2, "taskA2", "I am ContentA2.", uuidUA, "2020-12-08")
	taskB1 = entity.NewTask(uuidTB1, "taskB1", "I am ContentB1.", uuidUB, "2020-12-08")
	taskB2 = entity.NewTask(uuidTB2, "taskB2", "I am ContentB2.", uuidUB, "2020-12-08")
)

func TestTaskRepository_Create(t *testing.T) {

	db, task := prepareTaskT(t)

	tests := []struct {
		name         string
		task         *entity.Task
		wantTask     *entity.Task
		wantErr      error
		prepareTasks []*entity.Task
	}{
		{
			name:     "新しいユーザーのタスクを正しく作成できる",
			task:     entity.NewTask(uuidTB1, "taskB1", "I am ContentB1.", uuidUB, "2020-12-08"),
			wantTask: entity.NewTask(uuidTB1, "taskB1", "I am ContentB1.", uuidUB, "2020-12-08"),
			wantErr:  nil,
			prepareTasks: []*entity.Task{
				taskA1,
			},
		},
		{
			name:     "既存のユーザーのタスクを正しく作成できる",
			task:     entity.NewTask(uuidTA2, "taskA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask: entity.NewTask(uuidTA2, "taskA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantErr:  nil,
			prepareTasks: []*entity.Task{
				taskA1,
			},
		},
		{
			name:         "Titleがnilの場合はErrMySQL",
			task:         entity.NewTask(uuidTA2, "", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'title' cannot be null"),
			prepareTasks: nil,
		},
		{
			name:         "Contentが空でも正しく作成できる",
			task:         entity.NewTask(uuidTA2, "taskA2", "", uuidUA, "2020-12-08"),
			wantTask:     entity.NewTask(uuidTA2, "taskA2", "", uuidUA, "2020-12-08"),
			wantErr:      nil,
			prepareTasks: nil,
		},
		{
			name:         "UserIDがnilの場合はErrMySQL",
			task:         entity.NewTask(uuidTA2, "tasksA2", "I am ContentA2.", "", "2020-12-08"),
			wantTask:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'user_id' cannot be null"),
			prepareTasks: nil,
		},
		{
			name:         "IsCompがtrueでも正しく作成できる",
			task:         entity.NewTask(uuidTB1, "taskB1", "I am ContentB1.", uuidUB, "2020-12-08").SetComp(true),
			wantTask:     entity.NewTask(uuidTB1, "taskB1", "I am ContentB1.", uuidUB, "2020-12-08").SetComp(true),
			wantErr:      nil,
			prepareTasks: nil,
		},
		{
			name:     "指定したIDのタスクが既に存在している場合はErrMySQL",
			task:     entity.NewTask(uuidTA2, "taskA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry '38397cad-8865-081f-3482-2a035f875d5c' for key 'tasks.PRIMARY'"),
			prepareTasks: []*entity.Task{
				taskA2,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareTaskTT(t, db, tt.prepareTasks)

			err := task.Create(db, tt.task)
			gotTask := tt.task

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotTask)
				return
			}
			if (tt.wantErr == nil) && (!taskEqual(t, gotTask, tt.wantTask)) {
				t.Errorf("Create() = %s, want %s", gotTask, tt.wantTask)
			}
		})
	}
	db.Exec("TRUNCATE TABLE tasks")
}

// addTaskData はテスト用のデータをデータベースに追加する
func addTaskData(t *testing.T, db *gorm.DB, tasks []*entity.Task) (err error) {
	t.Helper()
	for _, task := range tasks {
		err = db.Create(task).Error
		if err != nil {
			return
		}
	}
	return
}

// taskEqual はCreatedAt, UpdatedAt以外のTaskのフィールドが同じかどうか判定する
func taskEqual(t *testing.T, got *entity.Task, want *entity.Task) bool {
	t.Helper()
	return (got.ID == want.ID) &&
		(got.Title.Equals(want.Title)) &&
		(got.Content.Equals(want.Content)) &&
		(got.UserID.Equals(want.UserID)) &&
		(got.IsCompleted == want.IsCompleted) &&
		(got.Date.Equals(want.Date))
}

func prepareTaskT(t *testing.T) (db *gorm.DB, task *TaskRepository) {
	t.Helper()

	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	task = new(TaskRepository)
	db = dbRepo.Connect()
	// db.LogMode(true)

	return
}

func prepareTaskTT(t *testing.T, db *gorm.DB, tasks []*entity.Task) {
	t.Helper()

	// databaseを初期化する
	db.Exec("TRUNCATE TABLE tasks")

	// 事前データの準備
	err := addTaskData(t, db, tasks)
	if err != nil {
		t.Fatal(err)
	}
}