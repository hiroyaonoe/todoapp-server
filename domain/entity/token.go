package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Token は受け取った文字列をハッシュ化して扱うためのstruct
type Token struct {
	value        NullString
	is_encrypted bool
}

func NewToken(s string) Token {
	return Token{value: NewNullString(s), is_encrypted: false}
}

func (t *Token) String() string {
	return t.value.String()
}

func (t *Token) Set(str string) {
	t.value.Set(str)
}

func (t *Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *Token) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	t.Set(str)
	t.is_encrypted = false
	return nil
}

func (t *Token) IsNull() bool {
	return !t.value.IsNull()
}

func (t *Token) Equal(s *Token) bool {
	return t.value.Equal(s.value) && t.is_encrypted == s.is_encrypted
}

// Encrypt はトークンをハッシュ化する
func (t *Token) Encrypt() error {
	if t.is_encrypted {
		return fmt.Errorf("Already encrypted:%s", t)
	}
	if t.String() == "" {
		return nil
	}
	token := []byte(t.String())
	digest, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.Set(string(digest))
	t.is_encrypted = true

	return nil
}

// Authenticate は２つのトークンが同一のものか判定する
func (hashed *Token) Authenticate(plain *Token) error {
	if !hashed.is_encrypted || plain.is_encrypted {
		return fmt.Errorf("Invalid tokens")
	}
	p1 := []byte(hashed.String())
	p2 := []byte(plain.String())

	return bcrypt.CompareHashAndPassword(p1, p2)
}

// // BeforeSave はデータベースに保存する際のフック
// func (t *Token) BeforeSave(tx *gorm.DB) (err error) {
// 	return t.Encrypt()
// }

// Scan はデータベースの値をTokenにマッピングする
func (t *Token) Scan(value interface{}) error {
	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Invalid value:%s", value)
	}
	t.Set(string(str))
	t.is_encrypted = true
	return nil
}

// Value はTokenのフィールドのうちデータベースに保存するものを指定する
func (t Token) Value() (driver.Value, error) {
	return t.value.Value()
}
