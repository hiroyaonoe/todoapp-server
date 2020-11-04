/*
domain is Enterprise Business Rules.
どこにも依存しない．
*/
package domain

import (
	"time"
)

type User struct {
    ID int `gorm:"primary_key"`
    Name string
    Password string `sql:"not null"`
    Email *string `sql:"not null;unique"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// この struct はビジネスロジックだと思うので、 usecase で書くべきなのかと思ったけど、
// ここに定義した。
type UserForGet struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Email *string `json:"email"`
}

func (u *User) BuildForGet() UserForGet {
    user := UserForGet{}
    user.ID = u.ID
    user.Name = u.Name
    if u.Email != nil {
        user.Email = u.Email
    } else {
        empty := ""
        user.Email = &empty
    }
    return user
}