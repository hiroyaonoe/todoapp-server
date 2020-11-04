package infrastructure

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB struct {
    Host string
    Username string
    Password string
    DBName string
    Port string
    Connection *gorm.DB
}

func NewDB() *DB {
    c := NewConfig()
    return newDB(&DB{
        Host: c.DB.Production.Host,
        Username: c.DB.Production.Username,
        Password: c.DB.Production.Password,
        DBName: c.DB.Production.DBName,
        Port: c.DB.Production.Port,
    })
}

func NewTestDB() *DB {
    c := NewConfig()
    return newDB(&DB{
        Host: c.DB.Test.Host,
        Username: c.DB.Test.Username,
        Password: c.DB.Test.Password,
        DBName: c.DB.Test.DBName,
        Port: c.DB.Test.Port,
    })
}

func newDB(d *DB) *DB {
    db, err := gorm.Open("mysql", d.Username + ":" + d.Password + "@tcp(" + d.Host + d.Port + ")/" + d.DBName + "?charset=utf8&parseTime=True&loc=Local")
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