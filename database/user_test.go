package database

import (
	"reflect"
	"testing"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

const (
	uuidUA = "df1ecfbf-e5f8-5eab-d49c-3a3f2e201fa3"
	uuidUB = "38d94bb3-d13d-76c8-b4aa-54985158d899"
	uuidUZ = "98457fea-708f-bb8e-3e5e-fe1b43f1acad"
)

var (
	userA = entity.NewUser(uuidUA, "userA", "encrypted_passwordA", "exampleA@example.com")
	userB = entity.NewUser(uuidUB, "userB", "encrypted_passwordB", "exampleB@example.com")
)

func TestUserRepository_FindByID(t *testing.T) {

	db, user := prepareUserT(t)

	tests := []struct {
		name         string
		userid       string
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:     "正しくユーザが取得できる",
			userid:   uuidUB,
			wantUser: entity.NewUser(uuidUB, "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
				userB,
			},
		},
		{
			name:     "存在しないユーザーの場合はErrRecordNotFound",
			userid:   uuidUZ,
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, db, tt.prepareUsers)

			gotUser, err := user.FindByID(db, tt.userid)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("FindByID() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("FindByID() got = %s", gotUser)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("FindByID() got = %s, want %s", gotUser, tt.wantUser)
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

func TestUserRepository_Create(t *testing.T) {

	db, user := prepareUserT(t)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:     "正しくユーザを作成できる",
			user:     entity.NewUser(uuidUB, "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantUser: entity.NewUser(uuidUB, "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:         "Nameがnilの場合はErrMySQL",
			user:         entity.NewUser(uuidUB, "", "encrypted_passwordB", "exampleB@example.com"),
			wantUser:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'name' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:         "Passwordがnilの場合はErrMySQL",
			user:         entity.NewUser(uuidUB, "userB", "", "exampleB@example.com"),
			wantUser:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'password' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:         "Emailがnilの場合はErrMySQL",
			user:         entity.NewUser(uuidUB, "userB", "encrypted_passwordB", ""),
			wantUser:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'email' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:     "IDがnilの場合はErrMySQL",
			user:     entity.NewUser("", "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x554, "Field 'id' doesn't have a default value"),
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したIDのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser(uuidUA, "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry 'df1ecfbf-e5f8-5eab-d49c-3a3f2e201fa3' for key 'users.PRIMARY'"),
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したEmailのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser(uuidUA, "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry 'exampleB@example.com' for key 'users.email'"),
			prepareUsers: []*entity.User{
				userB,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, db, tt.prepareUsers)

			err := user.Create(db, tt.user)
			gotUser := tt.user

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotUser)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("Create() = %s, want %s", gotUser, tt.wantUser)
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

func TestUserRepository_Update(t *testing.T) {

	db, user := prepareUserT(t)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:     "全フィールドを変更できる",
			user:     entity.NewUser(uuidUA, "userAA", "encrypted_passwordAA", "exampleAA@example.com"),
			wantUser: entity.NewUser(uuidUA, "userAA", "encrypted_passwordAA", "exampleAA@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "Nameのみを変更できる",
			user:     entity.NewUser(uuidUA, "userAA", "", ""),
			wantUser: entity.NewUser(uuidUA, "userAA", "encrypted_passwordA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "Passwordのみを変更できる",
			user:     entity.NewUser(uuidUA, "", "encrypted_passwordAA", ""),
			wantUser: entity.NewUser(uuidUA, "userA", "encrypted_passwordAA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "Emailのみを変更できる",
			user:     entity.NewUser(uuidUA, "", "", "exampleAA@example.com"),
			wantUser: entity.NewUser(uuidUA, "userA", "encrypted_passwordA", "exampleAA@example.com"),
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "全フィールドがもとと同じでも実行できる",
			user:     entity.NewUser(uuidUA, "userA", "encrypted_passwordA", "exampleA@example.com"),
			wantUser: entity.NewUser(uuidUA, "userA", "encrypted_passwordA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "全フィールドが空でも実行できる",
			user:     entity.NewUser(uuidUA, "", "", ""),
			wantUser: entity.NewUser(uuidUA, "userA", "encrypted_passwordA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "IDが指定されていない場合はErrRecordNotFound",
			user:     entity.NewUser("", "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したIDのユーザーが存在しない場合はErrRecordNotFound",
			user:     entity.NewUser(uuidUZ, "userB", "encrypted_passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したEmailのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser(uuidUA, "", "", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry 'exampleB@example.com' for key 'users.email'"),
			prepareUsers: []*entity.User{
				userA,
				userB,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, db, tt.prepareUsers)

			err := user.Update(db, tt.user)
			gotUser := tt.user

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotUser)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("Create() = %s, want %s", gotUser, tt.wantUser)
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

func TestUserRepository_Delete(t *testing.T) {

	db, user := prepareUserT(t)

	tests := []struct {
		name         string
		userid       string
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:    "正しくユーザを削除できる",
			userid:  uuidUB,
			wantErr: nil,
			prepareUsers: []*entity.User{
				userA,
				userB,
			},
		},
		{
			name:    "存在しないユーザーの場合はErrRecordNotFound",
			userid:  uuidUZ,
			wantErr: entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, db, tt.prepareUsers)

			err := user.Delete(db, tt.userid)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Delete() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

// addUserData はテスト用のデータをデータベースに追加する
func addUserData(t *testing.T, db *gorm.DB, users []*entity.User) (err error) {
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
func userEqual(t *testing.T, got *entity.User, want *entity.User) bool {
	t.Helper()
	return (got.ID == want.ID) &&
		(got.Name.Equals(want.Name)) &&
		(got.Password.Equals(want.Password)) &&
		(got.Email.Equals(want.Email))
}

func prepareUserT(t *testing.T) (db *gorm.DB, user *UserRepository) {
	t.Helper()

	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user = new(UserRepository)
	db = dbRepo.Connect()
	// db.LogMode(true)

	return
}

func prepareUserTT(t *testing.T, db *gorm.DB, users []*entity.User) {
	t.Helper()

	// databaseを初期化する
	db.Exec("TRUNCATE TABLE users")

	// 事前データの準備
	err := addUserData(t, db, users)
	if err != nil {
		t.Fatal(err)
	}
}
