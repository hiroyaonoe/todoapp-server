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

// DB はデータベースの情報を示す TODO: フィールドをprivate化
type DB struct {
	Host       string
	Username   string
	Password   string
	DBName     string
	Port       string
	Connection *gorm.DB
}

func NewDB() *DB {
	c := config.NewConfig()
	return newDB(&DB{
		Host:     c.DB.Production.Host,
		Username: c.DB.Production.Username,
		Password: c.DB.Production.Password,
		DBName:   c.DB.Production.DBName,
		Port:     c.DB.Production.Port,
	})
}

func NewTestDB() *DB {
	c := config.NewConfig()
	return newDB(&DB{
		Host:     c.DB.Test.Host,
		Username: c.DB.Test.Username,
		Password: c.DB.Test.Password,
		DBName:   c.DB.Test.DBName,
		Port:     c.DB.Test.Port,
	})
}

func newDB(d *DB) *DB {
	db, err := gorm.Open("mysql", d.Username+":"+d.Password+"@tcp("+d.Host+d.Port+")/"+d.DBName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	d.Connection = db
	return d
}

func (db *DB) Begin() *gorm.DB {
	return db.Connection.Begin()
}

func (db *DB) Connect() *gorm.DB {
	return db.Connection
}

func (db *DB) Migrate() {
	connection := db.Connect()
	connection.AutoMigrate(&entity.User{})
	connection.AutoMigrate(&entity.Task{})
}
