package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/todoapp-server/interfaces/controllers"
)

type Routing struct {
    DB *DB
    Gin *gin.Engine
    Port string
}

func NewRouting(db *DB) *Routing {
    c := NewConfig()
    r := &Routing{
        DB: db,
        Gin: gin.Default(),
        Port: c.Routing.Port,
    }
    r.setRouting()
    return r
}

func (r *Routing) setRouting() {
    userController := controllers.NewUserController(r.DB)
    r.Gin.GET("/user/:id", func (c *gin.Context) { userController.Get(c) })
}

func (r *Routing) Run() {
    r.Gin.Run(r.Port)
}