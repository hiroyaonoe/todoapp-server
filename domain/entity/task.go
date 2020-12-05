package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Task は内部で処理する際のTask情報である
type Task struct {
	ID          NullString `gorm:"primary_key"`
	Title       NullString `gorm:"not null"`
	Content     NullString
	UserID      NullString `gorm:"not null;index"`
	IsCompleted bool
	Date        Date
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask is the constructor of Task.(値が""の場合はsql.NullStringのnullとして扱う)
func NewTask(id string, title string, content string, userid string, date Date) (u *Task) {
	u = &Task{
		ID:          NewNullString(id),
		Title:       NewNullString(title),
		Content:     NewNullString(content),
		UserID:      NewNullString(userid),
		IsCompleted: false,
		Date:        date,
	}
	return
}

// NewID はTaskのUUIDを生成
func (t *Task) NewID() *Task {
	id := uuid.New().String()
	t.ID = NewNullString(id)
	return t
}

// SetComp はTaskのIsCompletedを設定する
func (t *Task) SetComp(comp bool) *Task {
	t.IsCompleted = comp
	return t
}

// TaskForJSON はJSONにして外部に公開するTask情報である
type TaskForJSON struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"iscomp"`
	Date        string `json:"date"`
}

// ToTaskForJSON はTaskからTaskForJSONを取得する関数である
func (t *Task) ToTaskForJSON() (p *TaskForJSON) {
	p = &TaskForJSON{
		ID:          t.ID.String,
		Title:       t.Title.String,
		Content:     t.Content.String,
		IsCompleted: t.IsCompleted,
		Date:        t.Date.String(),
	}
	return
}

func (t *Task) String() (str string) {
	str = fmt.Sprintf("&entity.Task{ID:%s, Title:%s, Content:%s, UserID:%s, IsCompleted:%t, Date:%s, CreatedAt:%s, UpdatedAt: %s",
		t.ID.String, t.Title.String, t.Content.String, t.UserID.String, t.IsCompleted, t.Date, t.CreatedAt, t.UpdatedAt)
	return
}
