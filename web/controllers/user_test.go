package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/hiroyaonoe/todoapp-server/domain/mock_repository"
)

const (
	uuidUA = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
)

type testInfo struct {
	name                string            // test名
	userid              string            // cookieに入れるuserid
	params              map[string]string // context.Param
	body                string            // request body
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
			userid: uuidUA,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().FindByID(uuidUA).Return(&entity.User{
					ID:        entity.NewNullString(uuidUA),
					Name:      entity.NewNullString("username"),
					Password:  entity.NewToken("encrypted_password"),
					Email:     entity.NewNullString("example@example.com"),
					CreatedAt: time.Unix(100, 0),
					UpdatedAt: time.Unix(100, 0),
				}, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewUser(uuidUA, "username", "", "example@example.com"),
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: uuidUA,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().FindByID(uuidUA).Return(&entity.User{}, entity.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: ErrUserNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusUnauthorized",
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: ErrUnauthorized.Error(),
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
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Create(gomock.Any()).
					DoAndReturn(func(user *entity.User) error {
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
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name: "Requestにnameが含まれていないならStatusBadRequest",
			body: `{
				"password":"password",
				"email":"example@example.com"
				}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name: "Requestにpasswordが含まれていないならStatusBadRequest",
			body: `{
				"name":"username",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name: "Requestにemailが含まれていないならStatusBadRequest",
			body: `{
				"name":"username",
				"password":"password",
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name: "RequestBodyが不正ならStatusBadRequest",
			body: `{
				"id":10,
				"title":"title",
				"userid":3,
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name: "RequestBodyがJSONでないならStatusBadRequest",
			body: `aaaaa`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name: "同じemailのユーザーが既に存在するならばErrDuplicatedEmail",
			body: `{
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Create(gomock.Any()).Return(
					entity.NewErrMySQL(0x426, "Duplicate entry 'example@example.com' for key 'users.email'"))
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrDuplicatedEmail.Error(),
		},
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
			name:   "正しくフィールドを更新出来る",
			userid: uuidUA,
			body: `{
				"name":"newname",
				"password":"newpassword",
				"email":"newexample@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any()).
					DoAndReturn(func(user *entity.User) error {
						user.SetID(uuidUA)
						user.Password = entity.NewToken("encrypted_newpassword")
						user.CreatedAt = time.Unix(100, 0)
						user.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: entity.NewUser(uuidUA, "newname", "", "newexample@example.com"),
		},
		{
			name:   "nameがないならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name:   "passwordがないならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"name":"newname",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name:   "emailがないならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"name":"newname",
				"password":"password"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name:   "RequestBodyが不正ならStatusBadRequest",
			userid: uuidUA,
			body: `{
				"id":10,
				"title":"title",
				"userid":3
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name:   "RequestBodyがJSONでないならStatusBadRequest",
			userid: uuidUA,
			body:   `aaaaa`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrBadRequest.Error(),
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: uuidUA,
			body: `{
				"name":"newname",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any()).Return(entity.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: ErrUserNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusUnauthorized",
			body: `{
				"name":"newname",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: ErrUnauthorized.Error(),
		},
		{
			name:   "同じemailのユーザーが既に存在するならばErrDuplicatedEmail",
			userid: uuidUA,
			body: `{
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any()).Return(
					entity.NewErrMySQL(0x426, "Duplicate entry 'example@example.com' for key 'users.email'"))
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
			wantData: ErrDuplicatedEmail.Error(),
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
			userid: uuidUA,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Delete(uuidUA).Return(nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
			wantData: nil,
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: uuidUA,
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Delete(uuidUA).Return(entity.ErrRecordNotFound)
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
			wantData: ErrUserNotFound.Error(),
		},
		{
			name: "Cookieが空ならStatusUnauthorized",
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
			wantData: ErrUnauthorized.Error(),
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
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	tt.prepareMockUserRepo(userRepo)

	userController = NewUserController(userRepo)
	return
}

func compareResult(t *testing.T, w *httptest.ResponseRecorder, tt testInfo) {
	t.Helper()

	if w.Code != tt.wantCode {
		t.Errorf("Code (-want +got) =\n- %d\n+ %d", tt.wantCode, w.Code)
	}

	var want, got map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	var jsonByte []byte
	if tt.wantErr {
		wantErrRes := newErrorRes(tt.wantCode, tt.wantData.(string))
		jsonByte, err = json.Marshal(wantErrRes)
	} else {
		jsonByte, err = json.Marshal(tt.wantData)
	}
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(jsonByte, &want)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Data (-want +got) =\n%s\n", diff)
	}
}
