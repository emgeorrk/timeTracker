package routes

import (
	"github.com/gin-gonic/gin"
	"timeTracker/app"
	"timeTracker/handlers"
)

func InitRouter(myApp *app.App) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", myApp.DB)
	})

	v1 := r.Group("/api/v1")

	v1.GET("/users", handlers.GetUsers)
	v1.GET("/users/:id/tasks_overview", handlers.GetTasksOverview)

	v1.POST("/users", handlers.CreateUser)
	v1.POST("/users/:id/tasks", handlers.CreateTask)
	v1.POST("/users/:id/tasks/:task_id/start", handlers.StartTaskTimer)
	v1.POST("/users/:id/tasks/:task_id/stop", handlers.StopTaskTimer)

	v1.PUT("/users/:id", handlers.UpdateUser)

	v1.DELETE("/users/:id", handlers.DeleteUser)

	return r
}
