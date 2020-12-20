package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/mock_repository"
	"github.com/jinzhu/gorm"
)

// user_test上にあるので不要
// const (
// 	uuid = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
// )
// func TestMain(m *testing.M) {
// 	gin.SetMode("test")
// 	m.Run()
// }

func TestTaskController_Create(t *testing.T) {

	tests := []testInfo{
		{
			name:   "正しくタスクを作成できる",
			userid: uuid,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(db *gorm.DB, task *entity.Task) error {
						task.SetID("any id")
						task.CreatedAt = time.Unix(100, 0)
						task.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewTask("any id", "taskname", "I am content.", "", "2020-12-06"),
		},
		{
			name:   "RequestにtaskIDが含まれているならStatusBadRequest",
			userid: uuid,
			body: `{
				"id":taskid,
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "Requestにtitleが含まれていないならStatusBadRequest",
			userid: uuid,
			body: `{
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "Requestにdateが含まれていないならStatusBadRequest",
			userid: uuid,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "dateのformatが不正ならStatusBadRequest",
			userid: uuid,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"invalid date"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "contentが含まれていなくてもok",
			userid: uuid,
			body: `{
				"title":"taskname",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(db *gorm.DB, task *entity.Task) error {
						task.SetID("any id")
						task.CreatedAt = time.Unix(100, 0)
						task.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantData: entity.NewTask("any id", "taskname", "", "", "2020-12-06"),
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name:   "iscompが含まれていなければfalseに設定",
			userid: uuid,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"date":"2020-12-06"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(db *gorm.DB, task *entity.Task) error {
						task.SetID("any id")
						task.CreatedAt = time.Unix(100, 0)
						task.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantData: entity.NewTask("any id", "taskname", "I am content.", "", "2020-12-06"),
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "Cookieが空ならStatusBadRequest",
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context, w := prepareTaskTT(t)

			// httpRequest
			context.Request, _ = http.NewRequest("POST", "/task", bytes.NewBufferString(tt.body))
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}

			// モック,コントローラーの準備
			ctrl, taskController := prepareMockTaskCtrl(t, tt)
			defer ctrl.Finish()

			taskController.Create(context)

			compareResult(t, w, tt)
		})
	}
}

func prepareTaskTT(t *testing.T) (context *gin.Context, w *httptest.ResponseRecorder) {
	t.Helper()
	t.Parallel()
	w = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(w)

	return
}

func prepareMockTaskCtrl(t *testing.T, tt testInfo) (ctrl *gomock.Controller, taskController *TaskController) {
	t.Helper()

	// モックの準備
	ctrl = gomock.NewController(t)
	dbRepo := mock_repository.NewMockDBRepository(ctrl)
	tt.prepareMockDBRepo(dbRepo)
	taskRepo := mock_repository.NewMockTaskRepository(ctrl)
	tt.prepareMockTaskRepo(taskRepo)

	taskController = NewTaskController(dbRepo, taskRepo)
	return
}
