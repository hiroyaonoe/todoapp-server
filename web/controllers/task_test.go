package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/mock_repository"
	"github.com/hiroyaonoe/todoapp-server/usecase"
	"github.com/jinzhu/gorm"
)

// user_test上にあるので不要
// const (
// 	uuid = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
// )

func TestTaskController_Create(t *testing.T) {

	tests := []struct {
		name                string
		userid              string
		body                string
		prepareMockDBRepo   func(db *mock_repository.MockDBRepository)
		prepareMockTaskRepo func(task *mock_repository.MockTaskRepository)
		wantData            interface{}
		wantErr             bool
		wantCode            int
	}{
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
			wantData: entity.TaskForJSON{
				ID:          "any id",
				Title:       "taskname",
				Content:     "I am content.",
				IsCompleted: false,
				Date:        entity.NewNullDate("2020-12-06"),
			},
			wantErr:  false,
			wantCode: http.StatusOK,
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
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
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
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
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
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
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
			wantData: entity.TaskForJSON{
				ID:          "any id",
				Title:       "taskname",
				IsCompleted: false,
				Date:        entity.NewNullDate("2020-12-06"),
			},
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
			wantData: entity.TaskForJSON{
				ID:          "any id",
				Title:       "taskname",
				IsCompleted: false,
				Date:        entity.NewNullDate("2020-12-06"),
			},
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
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode("test")
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)
			context.Request, _ = http.NewRequest("POST", "/task", bytes.NewBufferString(tt.body))
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}

			// モックの準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbRepo := mock_repository.NewMockDBRepository(ctrl)
			tt.prepareMockDBRepo(dbRepo)
			taskRepo := mock_repository.NewMockTaskRepository(ctrl)
			tt.prepareMockTaskRepo(taskRepo)

			taskController := &TaskController{
				Interactor: usecase.TaskInteractor{
					DB:   dbRepo,
					Task: taskRepo,
				},
			}

			taskController.Create(context)

			if w.Code != tt.wantCode {
				t.Errorf("Create() code = %d, want = %d", w.Code, tt.wantCode)
			}

			if tt.wantErr {
				actualData := ErrorForJSON{}
				expectData := tt.wantData.(ErrorForJSON)
				err := json.Unmarshal(w.Body.Bytes(), &actualData)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualData, expectData) {
					t.Errorf("Create() errData = %#v, want = %#v", actualData, expectData)
				}
			} else {
				actualData := entity.TaskForJSON{}
				expectData := tt.wantData.(entity.TaskForJSON)
				err := json.Unmarshal(w.Body.Bytes(), &actualData)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualData, expectData) {
					t.Errorf("Create() okData = %#v, want = %#v", actualData, expectData)
				}
			}
		})
	}
}
