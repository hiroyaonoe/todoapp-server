package main

import (
	"github.com/hiroyaonoe/todoapp-server/database"
	"github.com/hiroyaonoe/todoapp-server/web"
)

func main() {
	db := database.NewDB()
	// db := database.NewTestDB()
	db.Migrate()
	db.Connect().LogMode(true)
	user := database.NewUserRepository(db)
	task := new(database.TaskRepository)
	r := web.NewRouting(user, task)
	r.Run()
}
