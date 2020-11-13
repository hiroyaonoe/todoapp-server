package main

import (
	"github.com/hiroyaonoe/todoapp-server/database"
	"github.com/hiroyaonoe/todoapp-server/web"
)

func main() {
	db := database.NewDB()
	// db := database.NewTestDB()
	db.Migrate()
	user := new(database.UserRepository)
	r := web.NewRouting(db, user)
	r.Run()
}
