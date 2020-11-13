/*
Package web is Frameworks & Drivers.
ルーティング処理
*/
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/hiroyaonoe/todoapp-server/config"
	"github.com/hiroyaonoe/todoapp-server/database"
	"github.com/hiroyaonoe/todoapp-server/web/controllers"
)

type Routing struct {
	DB   *database.DB
	User *database.UserRepository
	Gin  *gin.Engine
	Port string
}

func NewRouting(db *database.DB, user *database.UserRepository) *Routing {
	c := config.NewConfig()
	r := &Routing{
		DB:   db,
		User: user,
		Gin:  gin.Default(),
		Port: c.Routing.Port,
	}
	r.setRouting()
	return r
}

func (r *Routing) setRouting() {
	// taskController := controllers.NewTaskController(r.DB)
	userController := controllers.NewUserController(r.DB, r.User)
	// authController := controllers.NewAuthController(r.DB)

	engine := r.Gin

	// middleware

	// engine.POST("/login", func(c *gin.Context) { authController.Login(c) })

	v1 := engine.Group("/api/v1")
	// tasks := v1.Group("/tasks")
	// tasks.POST("", func(c *gin.Context) { taskController.Create(c) })
	// tasks.GET("/:id", func(c *gin.Context) { taskController.GetbyID(c) })
	// tasks.PUT("/:id", func(c *gin.Context) { taskController.Update(c) })
	// tasks.Delete("/:id", func(c *gin.Context) { taskController.Delete(c) })
	// tasks.PUT("/:id/completed", func(c *gin.Context) { taskController.Switch(c) })
	// tasks.GET("/date/:date", func(c *gin.Context) { taskController.GetbyDate(c) })
	// tasks.GET("/date/from/:start/to/:end", func(c *gin.Context) { taskController.GetbyPeriod(c) })

	user := v1.Group("/user")
	user.GET("", func(c *gin.Context) { userController.Get(c) })
	user.POST("", func(c *gin.Context) { userController.Create(c) })
	user.PUT("", func(c *gin.Context) { userController.Update(c) })
	// user.DELETE("", func(c *gin.Context) { userController.Delete(c) })

}

func (r *Routing) Run() {
	r.Gin.Run(r.Port)
}
