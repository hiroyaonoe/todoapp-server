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
	task := database.NewTaskRepository(db)
	r := web.NewRouting(user, task)
	r.Run()
}
