package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Task は内部で処理する際のTask情報である
type Task struct {
	ID          NullString `gorm:"primary_key" json:"id"`
	Title       NullString `gorm:"not null" json:"title"`
	Content     NullString `json:"content"`
	UserID      NullString `gorm:"not null;index"`
	IsCompleted bool       `gorm:"not null" json:"iscomp"`
	Deadline    NullDate   `gorm:"not null" json:"deadline"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
}

// MarshalJSON はjsonにエンコードするときにUserIDフィールドを隠す
func (t *Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          NullString `json:"id"`
		Title       NullString `json:"title"`
		Content     NullString `json:"content"`
		IsCompleted bool       `json:"iscomp"`
		Deadline    NullDate   `json:"deadline"`
	}{
		ID:          t.ID,
		Title:       t.Title,
		Content:     t.Content,
		IsCompleted: t.IsCompleted,
		Deadline:    t.Deadline,
	})
}

// UnmarshalJSON はjsonからデコードするときにUserIDフィールドを隠さない
// func (t *Task) UnmarshalJSON(data []byte) error {
// }

// NewTask is the constructor of Task.(値が""の場合はsql.NullStringのnullとして扱う)
func NewTask(id string, title string, content string, userid string, deadline string) (u *Task) {
	u = &Task{
		ID:          NewNullString(id),
		Title:       NewNullString(title),
		Content:     NewNullString(content),
		UserID:      NewNullString(userid),
		IsCompleted: false,
		Deadline:    NewNullDate(deadline),
	}
	return
}

// NewID はTaskのUUIDを生成
func (t *Task) NewID() *Task {
	id := uuid.New().String()
	t.ID = NewNullString(id)
	return t
}

// SetID はTaskのIDを設定
func (t *Task) SetID(id string) *Task {
	t.ID = NewNullString(id)
	return t
}

// SetComp はTaskのIsCompletedを設定する
func (t *Task) SetComp(comp bool) *Task {
	t.IsCompleted = comp
	return t
}

// // TaskForJSON はJSONにして外部に公開するTask情報である
// type TaskForJSON struct {
// 	ID          string   `json:"id"`
// 	Title       string   `json:"title"`
// 	Content     string   `json:"content"`
// 	IsCompleted bool     `json:"iscomp"`
// 	Deadline    NullDate `json:"date"`
// }

// // ToTaskForJSON はTaskからTaskForJSONを取得する関数である
// func (t *Task) ToTaskForJSON() (p *TaskForJSON) {
// 	p = &TaskForJSON{
// 		ID:          t.ID.GetString(),
// 		Title:       t.Title.GetString(),
// 		Content:     t.Content.GetString(),
// 		IsCompleted: t.IsCompleted,
// 		Deadline:    t.Deadline,
// 	}
// 	return
// }

func (t *Task) String() (str string) {
	str = fmt.Sprintf("&entity.Task{ID:%s, Title:%s, Content:%s, UserID:%s, IsCompleted:%t, Deadline:%s, CreatedAt:%s, UpdatedAt: %s",
		t.ID.GetString(), t.Title.GetString(), t.Content.GetString(), t.UserID.GetString(), t.IsCompleted, t.Deadline, t.CreatedAt, t.UpdatedAt)
	return
}
