//go:generate mockgen -source=$GOFILE -destination=../mock_repository/mock_$GOFILE -package=mock_repository

package repository

import (
	"github.com/jinzhu/gorm"
)

type DBRepository interface {
	Begin() *gorm.DB
	Connect() *gorm.DB
}
