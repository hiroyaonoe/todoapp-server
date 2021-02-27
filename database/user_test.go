package database

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
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

			addUserData(t, user, tt.prepareUsers)

			gotUser, err := user.FindByID(tt.userid)

			if errorCompare(t, err, tt.wantErr) {
				t.Errorf("Data got = %s", gotUser)
			}
			if tt.wantErr == nil {
				cmpopt := cmpopts.IgnoreFields(entity.User{},
					"Password",
					"CreatedAt",
					"UpdatedAt")
				diff := cmp.Diff(tt.wantUser, gotUser, cmpopt)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
				if !gotUser.Password.Authenticate(&tt.wantUser.Password) {
					t.Errorf("Password (-want +got) =\n- %s\n+ %s",
						tt.wantUser.Password,
						gotUser.Password)
				}

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

			addUserData(t, user, tt.prepareUsers)

			err := user.Create(tt.user)
			gotUser := tt.user

			if errorCompare(t, err, tt.wantErr) {
				t.Errorf("Data got = %s", gotUser)
			}
			if tt.wantErr == nil {
				// IDは一致する必要なし
				cmpopt := cmpopts.IgnoreFields(entity.User{},
					"ID",
					"Password",
					"CreatedAt",
					"UpdatedAt")
				diff := cmp.Diff(tt.wantUser, gotUser, cmpopt)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
				if !gotUser.Password.Authenticate(&tt.wantUser.Password) {
					t.Errorf("Password (-want +got) =\n- %s\n+ %s",
						tt.wantUser.Password,
						gotUser.Password)
				}

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
			name:     "Nameが空ならばErrMySQL",
			user:     entity.NewUser(uuidUA, "", "passwordAA", "exampleAA@example.com"),
			wantUser: nil,
			wantErr:  entity.NewErrMySQL(0x418, "Column 'name' cannot be null"),
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "Passwordが空ならばErrMySQL",
			user:     entity.NewUser(uuidUA, "userAA", "", "exampleAA@example.com"),
			wantUser: nil,
			wantErr:  entity.NewErrMySQL(0x418, "Column 'password' cannot be null"),
			prepareUsers: []entity.User{
				userA,
			},
		},
		{
			name:     "Emailが空ならばErrMySQL",
			user:     entity.NewUser(uuidUA, "userAA", "passwordAA", ""),
			wantUser: nil,
			wantErr:  entity.NewErrMySQL(0x418, "Column 'email' cannot be null"),
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
			user:     entity.NewUser(uuidUA, "userA", "passwordA", "exampleB@example.com"),
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

			addUserData(t, user, tt.prepareUsers)

			err := user.Update(tt.user)
			gotUser := tt.user

			if errorCompare(t, err, tt.wantErr) {
				t.Errorf("Data got = %s", gotUser)
			}
			if tt.wantErr == nil {
				cmpopt := cmpopts.IgnoreFields(entity.User{},
					"Password",
					"CreatedAt",
					"UpdatedAt")
				diff := cmp.Diff(tt.wantUser, gotUser, cmpopt)
				if diff != "" {
					t.Errorf("Data (-want +got) =\n%s\n", diff)
				}
				if !gotUser.Password.Authenticate(&tt.wantUser.Password) {
					t.Errorf("Password (-want +got) =\n- %s\n+ %s",
						tt.wantUser.Password,
						gotUser.Password)
				}

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

			addUserData(t, user, tt.prepareUsers)

			err := user.Delete(tt.userid)

			errorCompare(t, err, tt.wantErr)
		})
	}
}

// addUserData はテスト用のユーザーデータをデータベースに追加する
func addUserData(t *testing.T, repo *UserRepository, users []entity.User) {
	t.Helper()

	// databaseを初期化する
	db := repo.db
	err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	err = db.Exec("TRUNCATE TABLE users").Error
	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		t.Fatal(err)
	}

	for _, user := range users {
		user.EncryptPassword()
		err = db.Create(&user).Error
		if err != nil {
			t.Fatal(err)
		}
	}
	return
}

func prepareUserT(t *testing.T) (user *UserRepository) {
	t.Helper()

	// dbに接続
	db := NewTestDB()
	user = NewUserRepository(db)
	// db.LogMode(true)

	return
}
