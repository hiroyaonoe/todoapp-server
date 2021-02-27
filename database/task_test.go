package database

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
)

const (
	uuidTA1 = "cb8aa4fc-1964-a965-00ea-8b158c0ffcc7"
	uuidTA2 = "38397cad-8865-081f-3482-2a035f875d5c"
	uuidTB1 = "2265150f-a3c9-d21e-ee75-a8c42a07807e"
	uuidTB2 = "b67a753e-90c0-4435-ede7-fdcda1504ca5"
)

var (
	taskA1 = *entity.NewTask(uuidTA1, "taskA1", "I am ContentA1.", uuidUA, "2020-12-08")
	taskA2 = *entity.NewTask(uuidTA2, "taskA2", "I am ContentA2.", uuidUA, "2020-01-05")
	taskB1 = *entity.NewTask(uuidTB1, "taskB1", "I am ContentB1.", uuidUB, "2020-12-09")
	taskB2 = *entity.NewTask(uuidTB2, "taskB2", "I am ContentB2.", uuidUB, "2020-01-06")
)

func TestTaskRepository_Create(t *testing.T) {

	task := prepareTaskT(t)

	tests := []struct {
		name         string
		task         *entity.Task
		wantTask     *entity.Task
		wantErr      error
		prepareTasks []entity.Task
	}{
		{
			name:         "既存のユーザーのタスクを正しく作成できる",
			task:         entity.NewTask("", "taskA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask:     entity.NewTask("any id", "taskA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantErr:      nil,
			prepareTasks: []entity.Task{},
		},
		{
			name:         "存在しないユーザーのタスクを追加すればErrMySQL",
			task:         entity.NewTask("", "taskB1", "I am ContentB1.", uuidUZ, "2020-12-08"),
			wantTask:     nil,
			wantErr:      entity.NewErrMySQL(0x5ac, "Cannot add or update a child row: a foreign key constraint fails (`golang`.`tasks`, CONSTRAINT `tasks_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`))"),
			prepareTasks: []entity.Task{},
		},
		{
			name:         "Titleがnilの場合はErrMySQL",
			task:         entity.NewTask("", "", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask:     nil,
			wantErr:      entity.NewErrMySQL(0x418, "Column 'title' cannot be null"),
			prepareTasks: nil,
		},
		{
			name:         "Contentが空でも正しく作成できる",
			task:         entity.NewTask("", "taskA2", "", uuidUA, "2020-12-08"),
			wantTask:     entity.NewTask("any id", "taskA2", "", uuidUA, "2020-12-08"),
			wantErr:      nil,
			prepareTasks: nil,
		},
		{
			name:         "UserIDがnilの場合はErrMySQL",
			task:         entity.NewTask("", "tasksA2", "I am ContentA2.", "", "2020-12-08"),
			wantTask:     nil,
			wantErr:      entity.NewErrMySQL(0x418, "Column 'user_id' cannot be null"),
			prepareTasks: nil,
		},
		{
			name:         "IsCompがtrueでも正しく作成できる",
			task:         entity.NewTask("", "taskB1", "I am ContentB1.", uuidUB, "2020-12-08").SetComp(true),
			wantTask:     entity.NewTask("any id", "taskB1", "I am ContentB1.", uuidUB, "2020-12-08").SetComp(true),
			wantErr:      nil,
			prepareTasks: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			addTaskData(t, task, tt.prepareTasks)

			err := task.Create(tt.task)
			gotTask := tt.task

			if errorCompare(t, err, tt.wantErr) {
				t.Errorf("Data got = %s", gotTask)
			}
			if tt.wantErr == nil {
				// IDは一致する必要なし
				cmpopt := cmpopts.IgnoreFields(entity.Task{},
					"ID",
					"CreatedAt",
					"UpdatedAt")
				if diff := cmp.Diff(tt.wantTask, gotTask, cmpopt); diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestTaskRepository_FindByID(t *testing.T) {

	task := prepareTaskT(t)

	tests := []struct {
		name         string
		tid          string
		uid          string
		wantTask     *entity.Task
		wantErr      error
		prepareTasks []entity.Task
	}{
		{
			name:     "正しくTaskを取得できる",
			tid:      uuidTA1,
			uid:      uuidUA,
			wantTask: entity.NewTask(uuidTA1, "taskA1", "I am ContentA1.", uuidUA, "2020-12-08"),
			wantErr:  nil,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "tidが存在しないならErrRecordNotFound",
			tid:      uuidTA2,
			uid:      uuidUA,
			wantTask: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "tidが存在してもuidが存在しないならErrRecordNotFound",
			tid:      uuidTA1,
			uid:      uuidUB,
			wantTask: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			addTaskData(t, task, tt.prepareTasks)

			gotTask, err := task.FindByID(tt.tid, tt.uid)

			if errorCompare(t, err, tt.wantErr) {
				t.Errorf("Data got = %s", gotTask)
			}
			if tt.wantErr == nil {
				cmpopt := cmpopts.IgnoreFields(entity.Task{},
					"CreatedAt",
					"UpdatedAt")
				if diff := cmp.Diff(tt.wantTask, gotTask, cmpopt); diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestTaskRepository_Update(t *testing.T) {

	task := prepareTaskT(t)

	tests := []struct {
		name         string
		task         *entity.Task
		wantTask     *entity.Task
		wantErr      error
		prepareTasks []entity.Task
	}{
		{
			name:     "全フィールドを変更できる",
			task:     entity.NewTask(uuidTA1, "taskA2", "I am ContentA2.", uuidUA, "2020-01-05"),
			wantTask: entity.NewTask(uuidTA1, "taskA2", "I am ContentA2.", uuidUA, "2020-01-05"),
			wantErr:  nil,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "全フィールドがもとと同じでも変更できる",
			task:     entity.NewTask(uuidTA1, "taskA1", "I am ContentA1.", uuidUA, "2020-12-08"),
			wantTask: entity.NewTask(uuidTA1, "taskA1", "I am ContentA1.", uuidUA, "2020-12-08"),
			wantErr:  nil,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "Contentが空でも変更できる",
			task:     entity.NewTask(uuidTA1, "taskA2", "", uuidUA, "2020-01-05"),
			wantTask: entity.NewTask(uuidTA1, "taskA2", "", uuidUA, "2020-01-05"),
			wantErr:  nil,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "指定したIDのTaskが存在しない場合はErrRecordNotFound",
			task:     entity.NewTask(uuidTA2, "tasksA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "指定したUserIDのTaskが存在しない場合はErrRecordNotFound",
			task:     entity.NewTask(uuidTA1, "tasksA2", "I am ContentA2.", uuidUZ, "2020-12-08"),
			wantTask: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "IDが指定されていない場合はErrRecordNotFound",
			task:     entity.NewTask("", "tasksA2", "I am ContentA2.", uuidUA, "2020-12-08"),
			wantTask: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:     "UserIDがnilの場合はErrMySQL",
			task:     entity.NewTask(uuidTA1, "tasksA2", "I am ContentA2.", "", "2020-12-08"),
			wantTask: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			addTaskData(t, task, tt.prepareTasks)

			err := task.Update(tt.task)
			gotTask := tt.task

			if errorCompare(t, err, tt.wantErr) {
				t.Errorf("Data got = %s", gotTask)
			}
			if tt.wantErr == nil {
				cmpopt := cmpopts.IgnoreFields(entity.Task{},
					"CreatedAt",
					"UpdatedAt")
				if diff := cmp.Diff(tt.wantTask, gotTask, cmpopt); diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
			}
		})
	}
}

func TestTaskRepository_Delete(t *testing.T) {

	task := prepareTaskT(t)

	tests := []struct {
		name         string
		taskid       string
		userid       string
		wantErr      error
		prepareTasks []entity.Task
	}{
		{
			name:    "Taskを削除できる",
			taskid:  uuidTA1,
			userid:  uuidUA,
			wantErr: nil,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
		{
			name:         "Taskが存在しないならErrRecordNotFound",
			taskid:       uuidTA1,
			userid:       uuidUA,
			wantErr:      entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{},
		},
		{
			name:    "Taskが存在してもUserIDが異なるならErrRecordNotFound",
			taskid:  uuidTA1,
			userid:  uuidUB,
			wantErr: entity.ErrRecordNotFound,
			prepareTasks: []entity.Task{
				taskA1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			addTaskData(t, task, tt.prepareTasks)

			err := task.Delete(tt.taskid, tt.userid)

			errorCompare(t, err, tt.wantErr)
		})
	}
}

// addTaskData はテスト用のタスクデータをデータベースに追加する
func addTaskData(t *testing.T, repo *TaskRepository, tasks []entity.Task) {
	t.Helper()

	// databaseを初期化する
	db := repo.db
	err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	err = db.Exec("TRUNCATE TABLE tasks").Error
	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		t.Fatal(err)
	}

	for _, task := range tasks {
		err = db.Create(&task).Error
		if err != nil {
			t.Fatal(err)
		}
	}
	return
}

func prepareTaskT(t *testing.T) (task *TaskRepository) {
	t.Helper()

	// dbに接続
	db := NewTestDB()
	task = NewTaskRepository(db)
	// db.LogMode(true)

	// Userデータの準備
	user := NewUserRepository(db)
	users := []entity.User{userA, userB}
	addUserData(t, user, users)

	return
}

// errorCompare はgo-cmpを利用してerrorを比較
func errorCompare(t *testing.T, got, want error) bool {
	t.Helper()

	if got == nil && want == nil {
		return false
	}
	if got == nil {
		t.Errorf("Error got = nil, want = %s", want)
		return true
	}
	if want == nil {
		t.Errorf("Error want = nil, got = %s", got)
		return true
	}
	if diff := cmp.Diff(want.Error(), got.Error()); diff != "" {
		t.Errorf("Error (-want +got) =\n%s", diff)
		return true
	}
	return false
}
