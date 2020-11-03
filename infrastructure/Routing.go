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
    taskController := controllers.NewTaskController(r.DB)
    userController := controllers.NewUserController(r.DB)
    authController := controllers.NewAuthController(r.DB)

    engine := r.Gin
    
    // middleware


    
    engine.GET("/login", func (c *gin.Context) { authController.Login(c) })

    v1 := engine.Group("/api/v1")
    v1.POST("/tasks", func (c *gin.Context) { taskController.Create(c) })
    v1.GET("/tasks/:id", func (c *gin.Context) { taskController.GetbyID(c) })
    v1.PUT("/tasks/:id", func (c *gin.Context) { taskController.Update(c) })
    v1.Delete("/tasks/:id", func (c *gin.Context) { taskController.Delete(c) })
    v1.PUT("/tasks/:id/completed", func (c *gin.Context) { taskController.Switch(c) })
    v1.GET("/tasks/date/:date", func (c *gin.Context) { taskController.GetbyDate(c) })
    v1.GET("/tasks/date/from/:start/to/:end", func (c *gin.Context) { taskController.GetbyPeriod(c) })

    v1.GET("/user", func (c *gin.Context) { userController.Get(c) })
    v1.POST("/user", func (c *gin.Context) { userController.Create(c) })
    v1.PUT("/user", func (c *gin.Context) { userController.Update(c) })
    v1.DELETE("/user", func (c *gin.Context) { userController.Delete(c) })




}

func (r *Routing) Run() {
    r.Gin.Run(r.Port)
}