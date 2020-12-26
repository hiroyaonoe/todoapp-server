package main

import (
	"github.com/hiroyaonoe/todoapp-server/database"
	"github.com/hiroyaonoe/todoapp-server/web"
)

func main() {
	db := database.NewDB()
	db.Migrate()
	db.Connect().LogMode(true)
	user := new(database.UserRepository)
	task := new(database.TaskRepository)
	r := web.NewRouting(db, user, task)
	r.Run()
}
