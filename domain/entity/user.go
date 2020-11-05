/*
entity is Enterprise Business Rules.
エンティティがどうあるべきか記述
どこにも依存しない．
*/
package entity

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Password  string    `sql:"not null" json:"-"`
	Email     string    `sql:"not null;unique" json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// この struct はビジネスロジックだと思うので、 usecase で書くべきなのかと思ったけど、
// ここに定義した。
type UserForGet struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

func (u *User) BuildForGet() User {
	user := User{}
	user.ID = u.ID
	user.Name = u.Name
	user.Email = u.Email
	return user
}