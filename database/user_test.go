package database

import (
	"testing"
	"time"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

var (
	userA = &entity.User{
				Name:"userA",
				Password:"passwordA",
				Email:"exampleA@example.com",
			}
	userB = &entity.User{
				Name:"userB",
				Password:"passwordB",
				Email:"exampleB@example.com",
			}
	userC = &entity.User{
				Name:"userC",
				Password:"passwordC",
				Email:"exampleC@example.com",
			}
)

func TestUserRepository_FindByID(t *testing.T) {
	tests := []struct {
		name     string
		userid   int
		wantUser *entity.User
		wantErr  error
		prepareUsers []*entity.User
	}{
		{
			name:"正しくユーザが取得できる",
			userid:3,
			wantUser: &entity.User{
				ID:3,
				Name:"userC",
				Password:"passwordC",
				Email:"exampleC@example.com",
				CreatedAt: time.Unix(100, 0),
				UpdatedAt: time.Unix(100, 0),
			},
			wantErr:nil,
			prepareUsers: []*entity.User{
				userA,
				userB,
				userC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbRepo := NewTestDB()
			dbRepo.Migrate()
			user := new(UserRepository)
			db := dbRepo.Connect()

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
			
			// databaseを初期化する
			db.Exec("TRUNCATE TABLE users")
		})
	}
}


// addData はテスト用のデータをデータベースに追加する
func addData(t *testing.T, db *gorm.DB, users []*entity.User) (err error){
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
