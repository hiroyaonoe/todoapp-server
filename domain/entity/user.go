/*
Package entity is Enterprise Business Rules.
エンティティがどうあるべきか記述
どこにも依存しない．
*/
package entity

import (
	"time"
)

// User は内部で処理する際のUser情報である
type User struct {
	ID        int `gorm:"primary_key"`
	Name      string
	Password  string `sql:"not null"`
	Email     string `sql:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserForJSON はJSONにして外部に公開するUser情報である
type UserForJSON struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToUserForJSON はUserからUserForJSONを取得する関数である
func (u *User) ToUserForJSON() (p *UserForJSON) {
	p = &UserForJSON{
		ID:u.ID,
		Name:u.Name,
		Email:u.Email,
	}
	return
}

// func (u *User) HidePassword() {
// 	u.Password = ""
// 	return
// }
