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
	"github.com/jinzhu/gorm"
)

const (
	uuid = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
)

type testInfo struct {
	name                string
	prepareMockDBRepo   func(db *mock_repository.MockDBRepository)
	prepareMockUserRepo func(user *mock_repository.MockUserRepository)
	wantErr             bool
	wantCode            int
	wantData            interface{}
}

func TestMain(m *testing.M) {
	gin.SetMode("test")
	m.Run()
}

func TestUserController_Get(t *testing.T) {

	tests := []struct {
		info   testInfo
		userid string
	}{
		{
			userid: uuid,
			info: testInfo{
				name: "正しくユーザが取得できる",
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
		},
		{
			userid: uuid,
			info: testInfo{
				name: "DBにユーザがいないときはErrUserNotFound",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
					db.EXPECT().Connect()
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
					user.EXPECT().FindByID(gomock.Any(), uuid).Return(&entity.User{}, entity.ErrRecordNotFound)
				},
				wantErr:  true,
				wantCode: http.StatusNotFound,
				wantData: entity.ErrUserNotFound.Error(),
			},
		},
		{
			info: testInfo{
				name: "Cookieが空ならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.info.name, func(t *testing.T) {
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
			ctrl, userController := prepareMockUserCtrl(t, tt.info)
			defer ctrl.Finish()

			userController.Get(context)

			compareResult(t, w, tt.info)
		})
	}
}

func TestUserController_Create(t *testing.T) {

	tests := []struct {
		info testInfo
		body string
	}{
		{
			body: `{
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			info: testInfo{
				name: "正しくユーザを作成できる",
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
		},
		{
			body: `{
				"id":98457fea-708f-bb8e-3e5e-fe1b43f1acad,
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			info: testInfo{
				name: "RequestにuserIDが含まれているならStatusBadRequest",

				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			body: `{
				"password":"password",
				"email":"example@example.com"
			}`,
			info: testInfo{
				name: "Requestにnameが含まれていないならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			body: `{
				"name":"username",
				"email":"example@example.com"
			}`,
			info: testInfo{
				name: "Requestにpasswordが含まれていないならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			body: `{
				"name":"username",
				"password":"password",
			}`,
			info: testInfo{
				name: "Requestにemailが含まれていないならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			body: `{
				"id":10,
				"title":"title",
				"userid":3,
			}`,
			info: testInfo{
				name: "RequestBodyが不正ならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			body: `aaaaa`,
			info: testInfo{
				name: "RequestBodyがJSONでないならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.info.name, func(t *testing.T) {
			context, w := prepareUserTT(t)

			// httpRequest
			context.Request, _ = http.NewRequest("POST", "/user", bytes.NewBufferString(tt.body))

			// モック,コントローラーの準備
			ctrl, userController := prepareMockUserCtrl(t, tt.info)
			defer ctrl.Finish()

			userController.Create(context)

			compareResult(t, w, tt.info)

		})
	}
}

func TestUserController_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		info   testInfo
		userid string
		body   string
	}{
		{
			userid: uuid,
			body: `{
				"name":"newname"
			}`,
			info: testInfo{
				name: "正しくnameを更新出来る",
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
		},
		{
			userid: uuid,
			body: `{
				"id":10,
				"title":"title",
				"userid":3
			}`,
			info: testInfo{
				name: "RequestBodyが不正ならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			userid: uuid,
			body:   `aaaaa`,
			info: testInfo{
				name: "RequestBodyがJSONでないならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
		{
			userid: uuid,
			body: `{
				"name":"newname"
			}`,
			info: testInfo{
				name: "DBにユーザがいないときはErrUserNotFound",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
					db.EXPECT().Connect()
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
					user.EXPECT().Update(gomock.Any(), gomock.Any()).Return(entity.ErrRecordNotFound)
				},
				wantErr:  true,
				wantCode: http.StatusNotFound,
				wantData: entity.ErrUserNotFound.Error(),
			},
		},
		{
			body: `{
				"name":"newname"
			}`,
			info: testInfo{
				name: "Cookieが空ならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.info.name, func(t *testing.T) {
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
			ctrl, userController := prepareMockUserCtrl(t, tt.info)
			defer ctrl.Finish()

			userController.Update(context)

			compareResult(t, w, tt.info)
		})
	}
}

func TestUserController_Delete(t *testing.T) {

	tests := []struct {
		info   testInfo
		userid string
	}{
		{
			userid: uuid,
			info: testInfo{
				name: "正しくユーザーを削除できる",
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
		},
		{
			userid: uuid,
			info: testInfo{
				name: "DBにユーザがいないときはErrUserNotFound",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
					db.EXPECT().Connect()
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
					user.EXPECT().Delete(gomock.Any(), uuid).Return(entity.ErrRecordNotFound)
				},
				wantErr:  true,
				wantCode: http.StatusNotFound,
				wantData: entity.ErrUserNotFound.Error(),
			},
		},
		{
			info: testInfo{
				name: "Cookieが空ならStatusBadRequest",
				prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				},
				prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				},
				wantErr:  true,
				wantCode: http.StatusBadRequest,
				wantData: entity.ErrBadRequest.Error(),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.info.name, func(t *testing.T) {
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
			ctrl, userController := prepareMockUserCtrl(t, tt.info)
			defer ctrl.Finish()

			userController.Delete(context)

			compareResult(t, w, tt.info)
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
	t.Logf("%s", w.Body.Bytes())
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
	acstr, err := json.Marshal(&actual)
	wastr, err := json.Marshal(&actual)
	t.Errorf("json actual = %s, want = %s", acstr, wastr)
	if !reflect.DeepEqual(actual, want) {
		t.Errorf("actual = %#v, want = %#v", actual, want)
	}
}
