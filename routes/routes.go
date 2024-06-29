package routes

import (
	"github.com/gin-gonic/gin"
	"timeTracker/handlers"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	
	v1 := r.Group("/")
	
	v1.GET("/users", handlers.GetUsers)
	v1.GET("/users/:id", handlers.GetUserByID)
	v1.POST("/users", handlers.CreateUser)
	v1.PUT("/users/:id", handlers.UpdateUser)
	v1.DELETE("/users/:id", handlers.DeleteUser)
	
	return r
}
