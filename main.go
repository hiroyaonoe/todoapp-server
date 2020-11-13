package main

import (
	"github.com/hiroyaonoe/todoapp-server/web"
	"github.com/hiroyaonoe/todoapp-server/database"
)

func main() {
	db := database.NewDB()
	// db := database.NewTestDB()
	db.Migrate()
	r := web.NewRouting(db)
	r.Run()
}
