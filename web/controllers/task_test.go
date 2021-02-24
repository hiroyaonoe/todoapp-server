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
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/hiroyaonoe/todoapp-server/domain/mock_repository"
)

// user_test上にあるので不要
// const (
// 	uuidUA = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
// )
// func TestMain(m *testing.M) {
// 	gin.SetMode("test")
// 	m.Run()
// }

const (
	uuidTA = "65b77c66-99f1-985a-74d1-caccf54cda73"
)

func TestTaskController_Create(t *testing.T) {

	tests := []testInfo{
		{
			name:   "正しくタスクを作成できる",
			userid: uuidUA,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Create(gomock.Any()).
					DoAndReturn(func(task *entity.Task) error {
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
			userid: uuidUA,
			body: `{
				"id":taskid,
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "Requestにtitleが含まれていないならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "Requestにdateが含まれていないならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false
			}`,
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "dateのformatが不正ならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"invalid date"
			}`,
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "contentが含まれていなくてもok",
			userid: uuidUA,
			body: `{
				"title":"taskname",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Create(gomock.Any()).
					DoAndReturn(func(task *entity.Task) error {
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
			userid: uuidUA,
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"date":"2020-12-06"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Create(gomock.Any()).
					DoAndReturn(func(task *entity.Task) error {
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
			name: "Cookieが空ならStatusUnauthorized",
			body: `{
				"title":"taskname",
				"content":"I am content.",
				"iscomp":false,
				"date":"2020-12-06"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: errs.ErrUnauthorized.Error(),
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

func TestTaskController_GetByID(t *testing.T) {

	tests := []testInfo{
		{
			name:   "正しくタスクを取得できる",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().FindByID(uuidTA, uuidUA).Return(&entity.Task{
					ID:          entity.NewNullString(uuidTA),
					Title:       entity.NewNullString("title"),
					Content:     entity.NewNullString("I am Content."),
					UserID:      entity.NewNullString(uuidUA),
					IsCompleted: true,
					Date:        entity.NewNullDate("2020-12-27"),
					CreatedAt:   time.Unix(100, 0),
					UpdatedAt:   time.Unix(100, 0),
				}, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewTask(uuidTA, "title", "I am Content.", uuidUA, "2020-12-27").SetComp(true),
		},
		{
			name:   "DBにTaskがないときはErrTaskNotFound",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().FindByID(uuidTA, uuidUA).Return(&entity.Task{}, errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrTaskNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusUnauthorized",
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: errs.ErrUnauthorized.Error(),
		},
		{
			name:   "paramが空ならStatusBadRequest",
			userid: uuidUA,
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
			setParams(t, tt, context)

			// モック,コントローラーの準備
			ctrl, taskController := prepareMockTaskCtrl(t, tt)
			defer ctrl.Finish()

			taskController.GetByID(context)

			compareResult(t, w, tt)
		})
	}
}

func TestTaskController_Update(t *testing.T) {

	tests := []testInfo{
		{
			name:   "全フィールドを更新できる",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			body: `{
				"title":"newtitle",
				"content":"I am new content.",
				"iscomp":true,
				"date":"2020-01-05"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Update(gomock.Any()).
					DoAndReturn(func(task *entity.Task) error {
						task.CreatedAt = time.Unix(100, 0)
						task.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewTask(uuidTA, "newtitle", "I am new content.", "", "2020-01-05").SetComp(true),
		},
		{
			name:   "フィールドが足りないならStatusBadRequest",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			body: `{
				"iscomp":true,
				"date":"2020-01-05"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "RequestBodyが不正ならStatusBadRequest",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			body: `{
				"id":"taskid",
				"name":"newname"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "RequestBodyがJSONでないならStatusBadRequest",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			body:   `aaaaa`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "DBにTaskがないときはErrTaskNotFound",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			body: `{
				"title":"newtitle",
				"content":"I am new content.",
				"iscomp":true,
				"date":"2020-01-05"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Update(gomock.Any()).Return(errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrTaskNotFound.Error(),
		},
		{
			name:   "DBにUserがないときはErrTaskNotFound",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			body: `{
				"title":"newtitle",
				"content":"I am new content.",
				"iscomp":true,
				"date":"2020-01-05"
			}`,
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Update(gomock.Any()).Return(errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrTaskNotFound.Error(),
		},
		{
			name:   "Cookieが空ならStatusUnauthorized",
			params: map[string]string{"id": uuidTA},
			body: `{
				"title":"newtitle",
				"content":"I am new content.",
				"iscomp":true,
				"date":"2020-01-05"
			}`,
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: errs.ErrUnauthorized.Error(),
		},
		{
			name:   "TaskIDが空ならErrBadRequest",
			userid: uuidUA,
			body: `{
				"title":"newtitle",
				"content":"I am new content.",
				"iscomp":true,
				"date":"2020-01-05"
			}`,
			prepareMockTaskRepo: func(user *mock_repository.MockTaskRepository) {
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
			context.Request, _ = http.NewRequest("PUT", "/task", bytes.NewBufferString(tt.body))
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}
			setParams(t, tt, context)

			// モック,コントローラーの準備
			ctrl, taskController := prepareMockTaskCtrl(t, tt)
			defer ctrl.Finish()

			taskController.Update(context)

			compareResult(t, w, tt)
		})
	}
}

func TestTaskController_Delete(t *testing.T) {

	tests := []testInfo{
		{
			name:   "正しくTaskを削除できる",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Delete(uuidTA, uuidUA).Return(nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: nil,
		},
		{
			name:   "DBにTaskがないときはErrTaskNotFound",
			userid: uuidUA,
			params: map[string]string{"id": uuidTA},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
				task.EXPECT().Delete(uuidTA, uuidUA).Return(errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrTaskNotFound.Error(),
		},
		{
			name:   "Cookieが空ならStatusUnauthorized",
			params: map[string]string{"id": uuidTA},
			prepareMockTaskRepo: func(task *mock_repository.MockTaskRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: errs.ErrUnauthorized.Error(),
		},
		{
			name:   "TaskIDが空ならErrBadRequest",
			userid: uuidUA,
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
			context.Request, _ = http.NewRequest("PUT", "/task", bytes.NewBufferString(tt.body))
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}
			setParams(t, tt, context)

			// モック,コントローラーの準備
			ctrl, taskController := prepareMockTaskCtrl(t, tt)
			defer ctrl.Finish()

			taskController.Delete(context)

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
	taskRepo := mock_repository.NewMockTaskRepository(ctrl)
	tt.prepareMockTaskRepo(taskRepo)

	taskController = NewTaskController(taskRepo)
	return
}

func setParams(t *testing.T, tt testInfo, c *gin.Context) {
	t.Helper()

	var params gin.Params
	for k, v := range tt.params {
		param := gin.Param{Key: k, Value: v}
		params = append(params, param)
	}

	c.Params = params
}
