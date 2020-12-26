//go:generate mockgen -source=$GOFILE -destination=../mock_repository/mock_$GOFILE -package=mock_repository

package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/hiroyaonoe/todoapp-server/domain/entity"
)

// TaskRepository is interface of Task
type TaskRepository interface {
	Create(db *gorm.DB, t *entity.Task) (err error)
	FindByID(db *gorm.DB, tid string, uid string) (task *entity.Task, err error)
	// 	FindByUser(db *gorm.DB, uid int) (tasks []entity.Task, err error)
	// 	Update(db *gorm.DB, t entity.Task) (task entity.Task, err error)
	// 	Delete(db *gorm.DB, id int) (taskid int, err error)
}
