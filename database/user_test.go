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
	userA = *entity.NewUser(uuidUA, "userA", "passwordA", "exampleA@example.com")
	userB = *entity.NewUser(uuidUB, "userB", "passwordB", "exampleB@example.com")
)

func TestUserRepository_FindByID(t *testing.T) {

	user := prepareUserT(t)

	tests := []struct {
		name         string
		userid       string
		wantUser     *entity.User
		wantErr      error
		prepareUsers []entity.User
	}{
		{
			name:     "正しくユーザが取得できる",
			userid:   uuidUB,
			wantUser: entity.NewUser(uuidUB, "userB", "passwordB", "exampleB@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
				userB,
			},
		},
		{
			name:     "存在しないユーザーの場合はErrRecordNotFound",
			userid:   uuidUZ,
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []entity.User{
				userA,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, user, tt.prepareUsers)

			gotUser, err := user.FindByID(tt.userid)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("FindByID() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("FindByID() got = %s", gotUser)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser, false)) {
				t.Errorf("Create() = %s, want %s", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUserRepository_Create(t *testing.T) {

	user := prepareUserT(t)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []entity.User
	}{
		{
			name:     "正しくユーザを作成できる",
			user:     entity.NewUser("", "userB", "passwordB", "exampleB@example.com"),
			wantUser: entity.NewUser("any id", "userB", "passwordB", "exampleB@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:         "Nameがnilの場合はErrMySQL",
			user:         entity.NewUser("", "", "passwordB", "exampleB@example.com"),
			wantUser:     nil,
			wantErr:      entity.NewErrMySQL(0x418, "Column 'name' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:         "Passwordがnilの場合はErrMySQL",
			user:         entity.NewUser("", "userB", "", "exampleB@example.com"),
			wantUser:     nil,
			wantErr:      entity.NewErrMySQL(0x418, "Column 'password' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:         "Emailがnilの場合はErrMySQL",
			user:         entity.NewUser("", "userB", "passwordB", ""),
			wantUser:     nil,
			wantErr:      entity.NewErrMySQL(0x418, "Column 'email' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:     "指定したEmailのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser("", "userB", "passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.NewErrMySQL(0x426, "Duplicate entry 'exampleB@example.com' for key 'users.email'"),
			prepareUsers: []entity.User{
				userB,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, user, tt.prepareUsers)

			err := user.Create(tt.user)
			gotUser := tt.user

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotUser)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser, true)) {
				t.Errorf("Create() = %s, want %s", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {

	user := prepareUserT(t)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []entity.User
	}{
		{
			name:     "全フィールドを変更できる",
			user:     entity.NewUser(uuidUA, "userAA", "passwordAA", "exampleAA@example.com"),
			wantUser: entity.NewUser(uuidUA, "userAA", "passwordAA", "exampleAA@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "Nameのみを変更できる",
			user:     entity.NewUser(uuidUA, "userAA", "", ""),
			wantUser: entity.NewUser(uuidUA, "userAA", "passwordA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "Passwordのみを変更できる",
			user:     entity.NewUser(uuidUA, "", "passwordAA", ""),
			wantUser: entity.NewUser(uuidUA, "userA", "passwordAA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "Emailのみを変更できる",
			user:     entity.NewUser(uuidUA, "", "", "exampleAA@example.com"),
			wantUser: entity.NewUser(uuidUA, "userA", "passwordA", "exampleAA@example.com"),
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "全フィールドがもとと同じでも実行できる",
			user:     entity.NewUser(uuidUA, "userA", "passwordA", "exampleA@example.com"),
			wantUser: entity.NewUser(uuidUA, "userA", "passwordA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "全フィールドが空でも実行できる",
			user:     entity.NewUser(uuidUA, "", "", ""),
			wantUser: entity.NewUser(uuidUA, "userA", "passwordA", "exampleA@example.com"),
			wantErr:  nil,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "IDが指定されていない場合はErrRecordNotFound",
			user:     entity.NewUser("", "userB", "passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "指定したIDのユーザーが存在しない場合はErrRecordNotFound",
			user:     entity.NewUser(uuidUZ, "userB", "passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "指定したEmailのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser(uuidUA, "", "", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.NewErrMySQL(0x426, "Duplicate entry 'exampleB@example.com' for key 'users.email'"),
			prepareUsers: []entity.User{
				userA,
				userB,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, user, tt.prepareUsers)

			err := user.Update(tt.user)
			gotUser := tt.user

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotUser)
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser, false)) {
				t.Errorf("Create() = %s, want %s", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {

	user := prepareUserT(t)

	tests := []struct {
		name         string
		userid       string
		wantErr      error
		prepareUsers []entity.User
	}{
		{
			name:    "正しくユーザを削除できる",
			userid:  uuidUB,
			wantErr: nil,
			prepareUsers: []entity.User{
				userA,
				userB,
			},
		},
		{
			name:    "存在しないユーザーの場合はErrRecordNotFound",
			userid:  uuidUZ,
			wantErr: entity.ErrRecordNotFound,
			prepareUsers: []entity.User{
				userA,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			prepareUserTT(t, user, tt.prepareUsers)

			err := user.Delete(tt.userid)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Delete() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
}

// addUserData はテスト用のデータをデータベースに追加する
func addUserData(t *testing.T, db *gorm.DB, users []entity.User) (err error) {
	t.Helper()
	for _, user := range users {
		user.EncryptPassword()
		err = db.Create(&user).Error
		if err != nil {
			return
		}
	}
	return
}

// userEqual はCreatedAt, UpdatedAt以外のUserのフィールドが同じかどうか判定する
func userEqual(t *testing.T, got *entity.User, want *entity.User, is_create bool) bool {
	t.Helper()
	ret := (got.Name.Equals(want.Name)) &&
		(got.Password.Authenticate(&want.Password) == nil) &&
		(got.Email.Equals(want.Email))
	if !is_create {
		ret = ret && (got.ID == want.ID)
	}
	return ret
}

func prepareUserT(t *testing.T) (user *UserRepository) {
	t.Helper()

	// dbに接続
	db := NewTestDB()
	db.Migrate()
	user = NewUserRepository(db)
	// db.LogMode(true)

	return
}

func prepareUserTT(t *testing.T, user *UserRepository, users []entity.User) {
	t.Helper()

	// databaseを初期化する
	user.db.Exec("TRUNCATE TABLE users")

	// 事前データの準備
	err := addUserData(t, user.db, users)
	if err != nil {
		t.Fatal(err)
	}
}
