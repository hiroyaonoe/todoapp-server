package main

import (
	"github.com/hiroyaonoe/todoapp-server/infrastructure"
)

func main() {
	db := infrastructure.NewDB()
	// db := infrastructure.NewTestDB()
	db.Migrate()
	r := infrastructure.NewRouting(db)
	r.Run()
}
