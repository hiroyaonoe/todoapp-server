/*
Package entity is Enterprise Business Rules.
エンティティがどうあるべきか記述
どこにも依存しない．
*/
package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User は内部で処理する際のUser情報である
type User struct {
	// ID        int        `gorm:"primary_key"`
	ID        NullString `gorm:"primary_key" json:"id"`
	Name      NullString `gorm:"not null" json:"name"`
	Password  NullString `gorm:"not null" json:"password"`
	Email     NullString `gorm:"not null;unique" json:"email"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

// MarshalJSON はjsonにエンコードするときにパスワードフィールドを隠す
func (u *User) MarshalJSON() ([]byte, error) {
	type alias User
	aliasUser := alias(*u)
	aliasUser.Password.Set("")

	return json.Marshal(&struct {
		alias
	}{
		aliasUser,
	})
}

// UnmarshalJSON はjsonからデコードするときにパスワードフィールドを隠さない
// func (u *User) UnmarshalJSON(data []byte) error {
// }

// NewUser is the constructor of User.(値が""の場合はsql.NullStringのnullとして扱う)
func NewUser(id string, name string, pass string, email string) (u *User) {
	u = &User{
		ID:       NewNullString(id),
		Name:     NewNullString(name),
		Password: NewNullString(pass),
		Email:    NewNullString(email),
	}
	return
}

// NewID はUserのUUIDを生成
func (u *User) NewID() *User {
	id := uuid.New().String()
	u.ID = NewNullString(id)
	return u
}

// SetID はUserのIDを設定
func (u *User) SetID(id string) *User {
	u.ID = NewNullString(id)
	return u
}

// // UserForJSON はJSONにして外部に公開するUser情報である
// type UserForJSON struct {
// 	ID    string `json:"id"`
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }

// // ToUserForJSON はUserからUserForJSONを取得する関数である
// func (u *User) ToUserForJSON() (p *UserForJSON) {
// 	p = &UserForJSON{
// 		ID:    u.ID.GetString(),
// 		Name:  u.Name.GetString(),
// 		Email: u.Email.GetString(),
// 	}
// 	return
// }

func (u *User) String() (str string) {
	str = fmt.Sprintf("&entity.User{ID:%s, Name:%s, Password:%s, Email:%s, CreatedAt:%s, UpdatedAt: %s",
		u.ID.GetString(), u.Name.GetString(), u.Password.GetString(), u.Email.GetString(), u.CreatedAt, u.UpdatedAt)
	return
}
