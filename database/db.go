/*
Package database is Frameworks & Drivers.
SQLへの接続，クエリはここで行う
domainにのみ依存
*/
package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiroyaonoe/todoapp-server/config"
	"github.com/hiroyaonoe/todoapp-server/domain/entity"
	"github.com/jinzhu/gorm"
)

// DB はデータベースの情報を示す
type DB struct {
	dsn        string
	connection *gorm.DB
}

func NewDB() *DB {
	return newDB(&DB{
		dsn: config.DSN("Production"),
	})
}

func NewTestDB() *DB {
	return newDB(&DB{
		dsn: config.DSN("Test"),
	})
}

func newDB(d *DB) *DB {
	db, err := gorm.Open("mysql", d.dsn)
	if err != nil {
		panic(err.Error())
	}
	d.connection = db
	return d
}

func (db *DB) Connect() *gorm.DB {
	return db.connection
}

func (db *DB) Migrate() {
	connection := db.Connect()
	connection.AutoMigrate(&entity.User{})
	connection.AutoMigrate(&entity.Task{})
}

func (db *DB) LogMode(b bool) {
	db.Connect().LogMode(b)
}
