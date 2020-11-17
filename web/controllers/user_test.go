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

func TestUserController_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		userid              string
		prepareMockDBRepo   func(db *mock_repository.MockDBRepository)
		prepareMockUserRepo func(user *mock_repository.MockUserRepository)
		wantData            interface{}
		wantErr             bool
		wantCode            int
	}{
		{
			name:   "正しくユーザが取得できる",
			userid: "3",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().FindByID(gomock.Any(), 3).Return(&entity.User{
					ID:        3,
					Name:      entity.NewNullString("username"),
					Password:  entity.NewNullString("password"),
					Email:     entity.NewNullString("example@example.com"),
					CreatedAt: time.Unix(100, 0),
					UpdatedAt: time.Unix(100, 0),
				}, nil)
			},
			wantData: entity.UserForJSON{
				ID:    3,
				Name:  "username",
				Email: "example@example.com",
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: "3",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().FindByID(gomock.Any(), 3).Return(&entity.User{}, entity.ErrRecordNotFound)
			},
			wantData: ErrorForJSON{
				Code: http.StatusNotFound,
				Err:  entity.ErrUserNotFound.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
		{
			name:   "Cookieが空ならStatusBadRequest",
			userid: "",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "CookieがintでないならStatusBadRequest",
			userid: "a",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
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
			context.Request, _ = http.NewRequest("GET", "/user", nil)
			context.Request.AddCookie(&http.Cookie{
				Name:  "id",
				Value: tt.userid,
			})

			// モックの準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbRepo := mock_repository.NewMockDBRepository(ctrl)
			tt.prepareMockDBRepo(dbRepo)
			userRepo := mock_repository.NewMockUserRepository(ctrl)
			tt.prepareMockUserRepo(userRepo)

			userController := &UserController{
				Interactor: usecase.UserInteractor{
					DB:   dbRepo,
					User: userRepo,
				},
			}

			userController.Get(context)

			if w.Code != tt.wantCode {
				t.Errorf("Get() code = %d, want = %d", w.Code, tt.wantCode)
			}

			if tt.wantErr {
				actualData := ErrorForJSON{}
				expectData := tt.wantData.(ErrorForJSON)
				err := json.Unmarshal(w.Body.Bytes(), &actualData)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualData, expectData) {
					t.Errorf("Get() errData = %#v, want = %#v", actualData, expectData)
				}
			} else {
				actualData := entity.UserForJSON{}
				expectData := tt.wantData.(entity.UserForJSON)
				err := json.Unmarshal(w.Body.Bytes(), &actualData)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualData, expectData) {
					t.Errorf("Get() okData = %#v, want = %#v", actualData, expectData)
				}
			}
		})
	}
}

func TestUserController_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                string
		prepareMockDBRepo   func(db *mock_repository.MockDBRepository)
		prepareMockUserRepo func(user *mock_repository.MockUserRepository)
		wantData            interface{}
		wantErr             bool
		wantCode            int
	}{
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
						user.ID = 3
						user.CreatedAt = time.Unix(100, 0)
						user.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantData: entity.UserForJSON{
				ID:    3,
				Name:  "username",
				Email: "example@example.com",
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "RequestにuserIDが含まれているならStatusBadRequest",
			body: `{
				"id":10,
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Requestにnameが含まれていないならStatusBadRequest",
			body: `{
				"id":10,
				"password":"password",
				"email":"example@example.com"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Requestにpasswordが含まれていないならStatusBadRequest",
			body: `{
				"id":10,
				"name":"username",
				"email":"example@example.com"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Requestにemailが含まれていないならStatusBadRequest",
			body: `{
				"id":10,
				"name":"username",
				"password":"password",
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "RequestBodyが不正ならStatusBadRequest",
			body: `{
				"id":10,
				"title":"title",
				"userid":3
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "RequestBodyがJSONでないならStatusBadRequest",
			body: `aaaaa`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
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
			context.Request, _ = http.NewRequest("POST", "/user", bytes.NewBufferString(tt.body))

			// モックの準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbRepo := mock_repository.NewMockDBRepository(ctrl)
			tt.prepareMockDBRepo(dbRepo)
			userRepo := mock_repository.NewMockUserRepository(ctrl)
			tt.prepareMockUserRepo(userRepo)

			userController := &UserController{
				Interactor: usecase.UserInteractor{
					DB:   dbRepo,
					User: userRepo,
				},
			}

			userController.Create(context)

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
				actualData := entity.UserForJSON{}
				expectData := tt.wantData.(entity.UserForJSON)
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

func TestUserController_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		userid              string
		body                string
		prepareMockDBRepo   func(db *mock_repository.MockDBRepository)
		prepareMockUserRepo func(user *mock_repository.MockUserRepository)
		wantData            interface{}
		wantErr             bool
		wantCode            int
	}{
		{
			name:   "正しくnameを更新出来る",
			userid: "3",
			body: `{
				"name":"newname"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(db *gorm.DB, user *entity.User) error {
						user.ID = 3
						user.Password = entity.NewNullString("password")
						user.Email = entity.NewNullString("example@example.com")
						user.CreatedAt = time.Unix(100, 0)
						user.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantData: entity.UserForJSON{
				ID:    3,
				Name:  "newname",
				Email: "example@example.com",
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name:   "RequestBodyが不正ならStatusBadRequest",
			userid: "3",
			body: `{
				"id":10,
				"title":"title",
				"userid":3
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantData: ErrorForJSON{
				Code: http.StatusBadRequest,
				Err:  entity.ErrBadRequest.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "DBにユーザがいないときはErrUserNotFound",
			userid: "3",
			body: `{
				"name":"newname"
			}`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
				db.EXPECT().Connect()
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
				user.EXPECT().Update(gomock.Any(), gomock.Any()).Return(entity.ErrRecordNotFound)
			},
			wantData: ErrorForJSON{
				Code: http.StatusNotFound,
				Err:  entity.ErrUserNotFound.Error(),
			},
			wantErr:  true,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gin.SetMode("test")
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)
			context.Request, _ = http.NewRequest("PUT", "/user", bytes.NewBufferString(tt.body))
			context.Request.AddCookie(&http.Cookie{
				Name:  "id",
				Value: tt.userid,
			})

			// モックの準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbRepo := mock_repository.NewMockDBRepository(ctrl)
			tt.prepareMockDBRepo(dbRepo)
			userRepo := mock_repository.NewMockUserRepository(ctrl)
			tt.prepareMockUserRepo(userRepo)

			userController := &UserController{
				Interactor: usecase.UserInteractor{
					DB:   dbRepo,
					User: userRepo,
				},
			}

			userController.Update(context)

			if w.Code != tt.wantCode {
				t.Errorf("Update() code = %d, want = %d", w.Code, tt.wantCode)
			}

			if tt.wantErr {
				actualData := ErrorForJSON{}
				expectData := tt.wantData.(ErrorForJSON)
				err := json.Unmarshal(w.Body.Bytes(), &actualData)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualData, expectData) {
					t.Errorf("Update() errData = %#v, want = %#v", actualData, expectData)
				}
			} else {
				actualData := entity.UserForJSON{}
				expectData := tt.wantData.(entity.UserForJSON)
				err := json.Unmarshal(w.Body.Bytes(), &actualData)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualData, expectData) {
					t.Errorf("Update() okData = %#v, want = %#v", actualData, expectData)
				}
			}
		})
	}
}
