//go:generate mockgen -source=$GOFILE -destination=../mock_repository/mock_$GOFILE -package=mock_repository

package repository

import (
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
)

// TaskRepository is interface of Task
type TaskRepository interface {
	Create(t *entity.Task) (err error)
	FindByID(tid string, uid string) (task *entity.Task, err error)
	// 	FindByUser(uid int) (tasks []*entity.Task, err error)
	Update(t *entity.Task) (err error)
	Delete(tid string, uid string) (err error)
}
