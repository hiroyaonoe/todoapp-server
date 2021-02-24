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
	host       string
	username   string
	password   string
	dbName     string
	port       string
	connection *gorm.DB
}

func NewDB() *DB {
	c := config.NewConfig()
	return newDB(&DB{
		host:     c.DB.Production.Host,
		username: c.DB.Production.Username,
		password: c.DB.Production.Password,
		dbName:   c.DB.Production.DBName,
		port:     c.DB.Production.Port,
	})
}

func NewTestDB() *DB {
	c := config.NewConfig()
	return newDB(&DB{
		host:     c.DB.Test.Host,
		username: c.DB.Test.Username,
		password: c.DB.Test.Password,
		dbName:   c.DB.Test.DBName,
		port:     c.DB.Test.Port,
	})
}

func newDB(d *DB) *DB {
	db, err := gorm.Open("mysql", d.username+":"+d.password+"@tcp("+d.host+d.port+")/"+d.dbName+"?charset=utf8&parseTime=True&loc=Local")
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
