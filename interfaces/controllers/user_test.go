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
		wantMessage         string
		wantData            entity.User
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
				user.EXPECT().FindByID(gomock.Any(), 3).Return(entity.User{
					ID:        3,
					Name:      "username",
					Password:  "password",
					Email:     "example@example.com",
					CreatedAt: time.Unix(100, 0),
					UpdatedAt: time.Unix(100, 0),
				}, nil)
			},
			wantMessage: "success",
			wantData: entity.User{
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
				user.EXPECT().FindByID(gomock.Any(), 3).Return(entity.User{}, entity.ErrRecordNotFound)
			},
			wantMessage: entity.ErrUserNotFound.Error(),
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusNotFound,
		},
		{
			name:   "Cookieが空ならStatusBadRequest",
			userid: "",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantMessage: "strconv.Atoi: parsing \"\": invalid syntax",
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusBadRequest,
		},
		{
			name:   "CookieがintでないならStatusBadRequest",
			userid: "a",
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantMessage: "strconv.Atoi: parsing \"a\": invalid syntax",
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			actualH := struct {
				Message string
				Data    entity.User
			}{}
			err := json.Unmarshal(w.Body.Bytes(), &actualH)
			if err != nil {
				t.Fatal(err)
			}
			if actualH.Message != tt.wantMessage {
				t.Errorf("Get() message = %s, want = %s", actualH.Message, tt.wantMessage)
			}
			if !reflect.DeepEqual(actualH.Data, tt.wantData) {
				t.Errorf("Get() user = %+v, want = %+v", actualH.Data, tt.wantData)
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
		wantMessage         string
		wantData            entity.User
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
			wantMessage: "success",
			wantData: entity.User{
				ID:    3,
				Name:  "username",
				Email: "example@example.com",
			},
			wantErr:  false,
			wantCode: http.StatusOK,
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
			wantMessage: entity.ErrInvalidUser.Error(),
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusBadRequest,
		},
		{
			name: "RequestBodyがJSONでないならStatusBadRequest",
			body: `aaaaa`,
			prepareMockDBRepo: func(db *mock_repository.MockDBRepository) {
			},
			prepareMockUserRepo: func(user *mock_repository.MockUserRepository) {
			},
			wantMessage: "invalid character 'a' looking for beginning of value",
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				t.Errorf("Get() code = %d, want = %d", w.Code, tt.wantCode)
			}

			actualH := struct {
				Message string
				Data    entity.User
			}{}
			err := json.Unmarshal(w.Body.Bytes(), &actualH)
			if err != nil {
				t.Fatal(err)
			}
			if actualH.Message != tt.wantMessage {
				t.Errorf("Get() message = %s, want = %s", actualH.Message, tt.wantMessage)
			}
			if !reflect.DeepEqual(actualH.Data, tt.wantData) {
				t.Errorf("Get() user = %+v, want = %+v", actualH.Data, tt.wantData)
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
		wantMessage         string
		wantData            entity.User
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
						user.Password = "password"
						user.Email = "example@example.com"
						user.CreatedAt = time.Unix(100, 0)
						user.UpdatedAt = time.Unix(100, 0)
						return nil
					})
			},
			wantMessage: "success",
			wantData: entity.User{
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
			wantMessage: entity.ErrInvalidUser.Error(),
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusBadRequest,
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
			wantMessage: entity.ErrUserNotFound.Error(),
			wantData:    entity.User{},
			wantErr:     true,
			wantCode:    http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				t.Errorf("Get() code = %d, want = %d", w.Code, tt.wantCode)
			}

			actualH := struct {
				Message string
				Data    entity.User
			}{}
			err := json.Unmarshal(w.Body.Bytes(), &actualH)
			if err != nil {
				t.Fatal(err)
			}
			if actualH.Message != tt.wantMessage {
				t.Errorf("Get() message = %s, want = %s", actualH.Message, tt.wantMessage)
			}
			if !reflect.DeepEqual(actualH.Data, tt.wantData) {
				t.Errorf("Get() user = %+v, want = %+v", actualH.Data, tt.wantData)
			}

		})
	}
}
