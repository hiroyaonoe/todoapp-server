package database

import (
	"testing"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

var (
	userA = &entity.User{
		Name:     "userA",
		Password: "passwordA",
		Email:    "exampleA@example.com",
	}
	userB = &entity.User{
		Name:     "userB",
		Password: "passwordB",
		Email:    "exampleB@example.com",
	}
	userC = &entity.User{
		Name:     "userC",
		Password: "passwordC",
		Email:    "exampleC@example.com",
	}
)

func TestUserRepository_FindByID(t *testing.T) {
	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user := new(UserRepository)
	db := dbRepo.Connect()
	db.LogMode(true)

	tests := []struct {
		name         string
		userid       int
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:   "正しくユーザが取得できる",
			userid: 2,
			wantUser: &entity.User{
				ID:       2,
				Name:     "userB",
				Password: "passwordB",
				Email:    "exampleB@example.com",
			},
			wantErr: nil,
			prepareUsers: []*entity.User{
				userA,
				userB,
			},
		},
		{
			name:     "存在しないユーザーの場合はErrRecordNotFound",
			userid:   2,
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// databaseを初期化する
			db.Exec("TRUNCATE TABLE users")

			// 事前データの準備
			// tx := db.Begin()
			// defer tx.Rollback()
			err := addData(t, db, tt.prepareUsers)
			if err != nil {
				t.Fatal(err)
			}

			gotUser, err := user.FindByID(db, tt.userid)

			if err != tt.wantErr {
				t.Errorf("FindByID() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("FindByID() = %#v, want %#v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user := new(UserRepository)
	db := dbRepo.Connect()
	db.LogMode(true)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name: "正しくユーザを作成できる",
			user: &entity.User{
				Name:     "userB",
				Password: "passwordB",
				Email:    "exampleB@example.com",
			},
			wantUser: &entity.User{
				ID:       2,
				Name:     "userB",
				Password: "passwordB",
				Email:    "exampleB@example.com",
			},
			wantErr: nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name: "Nameがnilの場合はErr",
			user: &entity.User{
				Password: "passwordD",
				Email:    "exampleD@example.com",
			},
			wantUser:     nil,
			wantErr:      entity.ErrRecordNotFound,
			prepareUsers: nil,
		},
		{
			name: "Passwordがnilの場合はErr",
			user: &entity.User{
				Name:  "userD",
				Email: "exampleD@example.com",
			},
			wantUser:     nil,
			wantErr:      entity.ErrRecordNotFound,
			prepareUsers: nil,
		},
		{
			name: "Emailがnilの場合はErr",
			user: &entity.User{
				Name:     "userD",
				Password: "passwordD",
			},
			wantUser:     nil,
			wantErr:      entity.ErrRecordNotFound,
			prepareUsers: nil,
		},
		{
			name: "IDが指定されている(0でない)場合はそのIDで作成",
			user: &entity.User{
				ID:       100,
				Name:     "userD",
				Password: "passwordD",
				Email:    "exampleD@example.com",
			},
			wantUser: &entity.User{
				ID:       100,
				Name:     "userD",
				Password: "passwordD",
				Email:    "exampleD@example.com",
			},
			wantErr: nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name: "指定したIDのユーザーが既に存在している場合はErr",
			user: &entity.User{
				ID:       100,
				Name:     "userD",
				Password: "passwordD",
				Email:    "exampleD@example.com",
			},
			wantUser: &entity.User{
				ID:       100,
				Name:     "userD",
				Password: "passwordD",
				Email:    "exampleD@example.com",
			},
			wantErr: entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// databaseを初期化する
			db.Exec("TRUNCATE TABLE users")

			// 事前データの準備
			err := addData(t, db, tt.prepareUsers)
			if err != nil {
				t.Fatal(err)
			}

			err = user.Create(db, tt.user)
			gotUser := tt.user

			if err != tt.wantErr {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("Create() = %#v, want %#v", gotUser, tt.wantUser)
			}
		})
	}
}

// addData はテスト用のデータをデータベースに追加する
func addData(t *testing.T, db *gorm.DB, users []*entity.User) (err error) {
	t.Helper()
	for _, user := range users {
		err = db.Create(user).Error
		if err != nil {
			return
		}
	}
	return
}

// userEqual はCreatedAt, UpdatedAt以外のUserのフィールドが同じかどうか判定する
func userEqual(t *testing.T, a *entity.User, b *entity.User) bool {
	t.Helper()
	return (a.ID == b.ID) &&
		(a.Name == b.Name) &&
		(a.Password == b.Password) &&
		(a.Email == b.Email)
}
