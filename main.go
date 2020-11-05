package main

import (
	"github.com/hiroyaonoe/todoapp-server/infrastructure"
)

func main() {
	db := infrastructure.NewDB()
	r := infrastructure.NewRouting(db)
	r.Run()
}
