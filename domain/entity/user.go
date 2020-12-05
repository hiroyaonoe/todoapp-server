/*
Package entity is Enterprise Business Rules.
エンティティがどうあるべきか記述
どこにも依存しない．
*/
package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User は内部で処理する際のUser情報である
type User struct {
	ID        int        `gorm:"primary_key"`
	UUID      uuid.UUID  `gorm:"not null;unique;index"`
	Name      NullString `gorm:"not null"`
	Password  NullString `gorm:"not null"`
	Email     NullString `gorm:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser is the constructor of User.(値が""の場合はsql.NullStringのnullとして扱う)
func NewUser(name string, pass string, email string) (u *User) {
	u = &User{
		Name:     NewNullString(name),
		Password: NewNullString(pass),
		Email:    NewNullString(email),
	}
	return
}

// AddID はUserのIDを設定
func (u *User) AddID(id int) *User {
	u.ID = id
	return u
}

// AddUUID はUserのUUIDを設定
func (u *User) AddUUID(uuid uuid.UUID) *User {
	u.UUID = uuid
	return u
}

// UserForJSON はJSONにして外部に公開するUser情報である
type UserForJSON struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// ToUserForJSON はUserからUserForJSONを取得する関数である
func (u *User) ToUserForJSON() (p *UserForJSON) {
	p = &UserForJSON{
		UUID:  u.UUID,
		Name:  u.Name.String,
		Email: u.Email.String,
	}
	return
}

func (u *User) String() (str string) {
	str = fmt.Sprintf("&entity.User{ID:%d, UUID:%s, Name:%s, Password:%s, Email:%s, CreatedAt:%s, UpdatedAt: %s",
		u.ID, u.UUID.String(), u.Name.String, u.Password.String, u.Email.String, u.CreatedAt, u.UpdatedAt)
	return
}
