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
	"github.com/hiroyaonoe/todoapp-server/domain/errs"
	"github.com/hiroyaonoe/todoapp-server/domain/mock_repository"
	"github.com/jinzhu/gorm"
)

const (
	uuid = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
)

type testInfo struct {
	name                string
	userid              string
	body                string
	prepareMockDBRepo   func(db *mock_repository.MockDBRepository)
	prepareMockUserRepo func(user *mock_repository.MockUserRepository)
	prepareMockTaskRepo func(task *mock_repository.MockTaskRepository)
	wantErr             bool
	wantCode            int
	wantData            interface{}
}

func TestMain(m *testing.M) {
	gin.SetMode("test")
	m.Run()
}

func TestUserController_Get(t *testing.T) {

	tests := []testInfo{
		{
			name:   "正しくユーザが取得できる",
			userid: uuid,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().FindByID(gomock.Any(), uuid).Return(&entity.User{
					ID:        entity.NewNullString(uuid),
					Name:      entity.NewNullString("username"),
					Password:  entity.NewNullString("encrypted_password"),
					Email:     entity.NewNullString("example@example.com"),
					CreatedAt: time.Unix(100, 0),
					UpdatedAt: time.Unix(100, 0),
				}, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewUser(uuid, "username", "", "example@example.com"),
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: uuid,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().FindByID(gomock.Any(), uuid).Return(&entity.User{}, errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrUserNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusBadRequest",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context, w := prepareUserTT(t)

			// httpRequest
			context.Request, _ = http.NewRequest("GET", "/user", nil)
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}

			// モック,コントローラーの準備
			ctrl, userController := prepareMockUserCtrl(t, tt)
			defer ctrl.Finish()

			userController.Get(context)

			compareResult(t, w, tt)
		})
	}
}

func TestUserController_Create(t *testing.T) {

	tests := []testInfo{
		{
			name: "正しくユーザを作成できる",
			body: `{
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(db *gorm.DB, user *entity.User) error {
						user.SetID("any id")
						user.CreatedAt = time.Unix(100, 0)
						user.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewUser("any id", "username", "", "example@example.com"),
		},
		{
			name: "RequestにuserIDが含まれているならStatusBadRequest",
			body: `{
				"id":98457fea-708f-bb8e-3e5e-fe1b43f1acad,
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name: "Requestにnameが含まれていないならStatusBadRequest",
			body: `{
				"password":"password",
				"email":"example@example.com"
				}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name: "Requestにpasswordが含まれていないならStatusBadRequest",
			body: `{
				"name":"username",
				"email":"example@example.com"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name: "Requestにemailが含まれていないならStatusBadRequest",
			body: `{
				"name":"username",
				"password":"password",
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name: "RequestBodyが不正ならStatusBadRequest",
			body: `{
				"id":10,
				"title":"title",
				"userid":3,
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name: "RequestBodyがJSONでないならStatusBadRequest",
			body: `aaaaa`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		// {
		// 	name: "同じemailのユーザーが既に存在するならば",
		// 	body: `{
		// 		"name":"username",
		// 		"password":"password",
		// 		"email":"example@example.com"
		// 	}`,
		// 	prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
		// 		db.EXPECT().Connect()
		// 	},
		// 	prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
		// 		user.EXPECT().Create(gomock.Any(), gomock.Any()).Return(
		// 			entity.NewErrMySQL(0x426, "Duplicate entry 'example@example.com' for key 'users.email'"))
		// 	},
		// 	wantErr:  false,
		// 	wantCode: http.StatusOK,
		// 	wantData: entity.NewUser("any id", "username", "", "example@example.com"),
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context, w := prepareUserTT(t)

			// httpRequest
			context.Request, _ = http.NewRequest("POST", "/user", bytes.NewBufferString(tt.body))

			// モック,コントローラーの準備
			ctrl, userController := prepareMockUserCtrl(t, tt)
			defer ctrl.Finish()

			userController.Create(context)

			compareResult(t, w, tt)

		})
	}
}

func TestUserController_Update(t *testing.T) {
	t.Parallel()

	tests := []testInfo{
		{
			name:   "正しくnameを更新出来る",
			userid: uuid,
			body: `{
				"name":"newname"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(db *gorm.DB, user *entity.User) error {
						user.SetID(uuid)
						user.Password = entity.NewNullString("encrypted_password")
						user.Email = entity.NewNullString("example@example.com")
						user.CreatedAt = time.Unix(100, 0)
						user.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewUser(uuid, "newname", "", "example@example.com"),
		},
		{
			name:   "RequestBodyが不正ならStatusBadRequest",
			userid: uuid,
			body: `{
				"id":10,
				"title":"title",
				"userid":3
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "RequestBodyがJSONでないならStatusBadRequest",
			userid: uuid,
			body:   `aaaaa`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: uuid,
			body: `{
				"name":"newname"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrUserNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusBadRequest",
			body: `{
				"name":"newname"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context, w := prepareUserTT(t)

			// httpRequest
			context.Request, _ = http.NewRequest("PUT", "/user", bytes.NewBufferString(tt.body))
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}

			// モック,コントローラーの準備
			ctrl, userController := prepareMockUserCtrl(t, tt)
			defer ctrl.Finish()

			userController.Update(context)

			compareResult(t, w, tt)
		})
	}
}

func TestUserController_Delete(t *testing.T) {

	tests := []testInfo{
		{
			name:   "正しくユーザーを削除できる",
			userid: uuid,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Delete(gomock.Any(), uuid).Return(nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: nil,
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: uuid,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Delete(gomock.Any(), uuid).Return(errs.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: errs.ErrUserNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusBadRequest",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: errs.ErrBadRequest.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context, w := prepareUserTT(t)

			// httpRequest
			context.Request, _ = http.NewRequest("DELETE", "/user", nil)
			if tt.userid != "" {
				context.Request.AddCookie(&http.Cookie{
					Name:  "id",
					Value: tt.userid,
				})
			}

			// モック,コントローラーの準備
			ctrl, userController := prepareMockUserCtrl(t, tt)
			defer ctrl.Finish()

			userController.Delete(context)

			compareResult(t, w, tt)
		})
	}
}

func prepareUserTT(t *testing.T) (context *gin.Context, w *httptest.ResponseRecorder) {
	t.Helper()
	t.Parallel()
	w = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(w)

	return
}

func prepareMockUserCtrl(t *testing.T, tt testInfo) (ctrl *gomock.Controller, userController *UserController) {
	t.Helper()

	// モックの準備
	ctrl = gomock.NewController(t)
	dbRepo := mock_repository.NewMockDBRepository(ctrl)
	tt.prepareMockDBRepo(dbRepo)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	tt.prepareMockUserRepo(userRepo)

	userController = NewUserController(dbRepo, userRepo)
	return
}

func compareResult(t *testing.T, w *httptest.ResponseRecorder, tt testInfo) {
	t.Helper()

	if w.Code != tt.wantCode {
		t.Errorf("code = %d, want = %d", w.Code, tt.wantCode)
	}

	var want, actual map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actual)

	var jsonByte []byte
	if tt.wantErr {
		wantErrRes := newErrorRes(tt.wantCode, tt.wantData.(string))
		jsonByte, err = json.Marshal(wantErrRes)
	} else {
		jsonByte, err = json.Marshal(tt.wantData)
	}
	err = json.Unmarshal(jsonByte, &want)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, want) {
		t.Errorf("actual = %#v, want = %#v", actual, want)
	}
}
