package database

import (
	"reflect"
	"testing"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

var (
	userA = entity.NewUser("userA", "passwordA", "exampleA@example.com")
	userB = entity.NewUser("userB", "passwordB", "exampleB@example.com")
	userC = entity.NewUser("userC", "passwordC", "exampleC@example.com")
)

func TestUserRepository_FindByID(t *testing.T) {
	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user := new(UserRepository)
	db := dbRepo.Connect()
	// db.LogMode(true)

	tests := []struct {
		name         string
		userid       int
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:     "正しくユーザが取得できる",
			userid:   2,
			wantUser: entity.NewUser("userB", "passwordB", "exampleB@example.com").AddID(2),
			wantErr:  nil,
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

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("FindByID() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("FindByID() got = %s", gotUser.ToString())
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("FindByID() got = %s, want %s", gotUser.ToString(), tt.wantUser.ToString())
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

func TestUserRepository_Create(t *testing.T) {
	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user := new(UserRepository)
	db := dbRepo.Connect()
	// db.LogMode(true)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:     "正しくユーザを作成できる",
			user:     entity.NewUser("userB", "passwordB", "exampleB@example.com"),
			wantUser: entity.NewUser("userB", "passwordB", "exampleB@example.com").AddID(2),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:         "Nameがnilの場合はErrMySQL",
			user:         entity.NewUser("", "passwordB", "exampleB@example.com"),
			wantUser:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'name' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:         "Passwordがnilの場合はErrMySQL",
			user:         entity.NewUser("userB", "", "exampleB@example.com"),
			wantUser:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'password' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:         "Emailがnilの場合はErrMySQL",
			user:         entity.NewUser("userB", "passwordB", ""),
			wantUser:     nil,
			wantErr:      entity.ErrMySQL(0x418, "Column 'email' cannot be null"),
			prepareUsers: nil,
		},
		{
			name:     "IDが指定されている(0でない)場合はそのIDで作成",
			user:     entity.NewUser("userB", "passwordB", "exampleB@example.com").AddID(100),
			wantUser: entity.NewUser("userB", "passwordB", "exampleB@example.com").AddID(100),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したIDのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser("userB", "passwordB", "exampleB@example.com").AddID(1),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry '1' for key 'users.PRIMARY'"),
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したEmailのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser("userB", "passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry 'exampleB@example.com' for key 'users.email'"),
			prepareUsers: []*entity.User{
				userB,
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

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotUser.ToString())
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("Create() = %s, want %s", gotUser.ToString(), tt.wantUser.ToString())
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

func TestUserRepository_Update(t *testing.T) {
	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user := new(UserRepository)
	db := dbRepo.Connect()
	// db.LogMode(true)

	tests := []struct {
		name         string
		user         *entity.User
		wantUser     *entity.User
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:     "全フィールドを変更できる",
			user:     entity.NewUser("userAA", "passwordAA", "exampleAA@example.com").AddID(1),
			wantUser: entity.NewUser("userAA", "passwordAA", "exampleAA@example.com").AddID(1),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "Nameのみを変更できる",
			user:     entity.NewUser("userAA", "", "").AddID(1),
			wantUser: entity.NewUser("userAA", "passwordA", "exampleA@example.com").AddID(1),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "Passwordのみを変更できる",
			user:     entity.NewUser("", "passwordAA", "").AddID(1),
			wantUser: entity.NewUser("userA", "passwordAA", "exampleA@example.com").AddID(1),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "Emailのみを変更できる",
			user:     entity.NewUser("", "", "exampleAA@example.com").AddID(1),
			wantUser: entity.NewUser("userA", "passwordA", "exampleAA@example.com").AddID(1),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "全フィールドがもとと同じでも実行できる",
			user:     entity.NewUser("userA", "passwordA", "exampleA@example.com").AddID(1),
			wantUser: entity.NewUser("userA", "passwordA", "exampleA@example.com").AddID(1),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "全フィールドが空でも実行できる",
			user:     entity.NewUser("", "", "").AddID(1),
			wantUser: entity.NewUser("userA", "passwordA", "exampleA@example.com").AddID(1),
			wantErr:  nil,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "IDが指定されていない場合はErrRecordNotFound",
			user:     entity.NewUser("userB", "passwordB", "exampleB@example.com"),
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したIDのユーザーが存在しない場合はErrRecordNotFound",
			user:     entity.NewUser("userB", "passwordB", "exampleB@example.com").AddID(100),
			wantUser: nil,
			wantErr:  entity.ErrRecordNotFound,
			prepareUsers: []*entity.User{
				userA,
			},
		},
		{
			name:     "指定したEmailのユーザーが既に存在している場合はErrMySQL",
			user:     entity.NewUser("", "", "exampleB@example.com").AddID(1),
			wantUser: nil,
			wantErr:  entity.ErrMySQL(0x426, "Duplicate entry 'exampleB@example.com' for key 'users.email'"),
			prepareUsers: []*entity.User{
				userA,
				userB,
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

			err = user.Update(db, tt.user)
			gotUser := tt.user

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Create() error = %#v, wantErr %#v", err, tt.wantErr)
				t.Errorf("Create() got = %s", gotUser.ToString())
				return
			}
			if (tt.wantErr == nil) && (!userEqual(t, gotUser, tt.wantUser)) {
				t.Errorf("Create() = %s, want %s", gotUser.ToString(), tt.wantUser.ToString())
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
}

func TestUserRepository_Delete(t *testing.T) {
	// dbに接続
	dbRepo := NewTestDB()
	dbRepo.Migrate()
	user := new(UserRepository)
	db := dbRepo.Connect()
	// db.LogMode(true)

	tests := []struct {
		name         string
		userid       int
		wantErr      error
		prepareUsers []*entity.User
	}{
		{
			name:    "正しくユーザを削除できる",
			userid:  2,
			wantErr: nil,
			prepareUsers: []*entity.User{
				userA,
				userB,
			},
		},
		{
			name:    "存在しないユーザーの場合はErrRecordNotFound",
			userid:  2,
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

			err = user.Delete(db, tt.userid)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Delete() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
	db.Exec("TRUNCATE TABLE users")
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
func userEqual(t *testing.T, got *entity.User, want *entity.User) bool {
	t.Helper()
	return (got.ID == want.ID) &&
		(got.Name.Equals(want.Name)) &&
		(got.Password.Equals(want.Password)) &&
		(got.Email.Equals(want.Email))
}
