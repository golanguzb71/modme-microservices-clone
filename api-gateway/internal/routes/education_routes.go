package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func EducationRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	room := api.Group("/room")
	{
		room.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateRoom)
		room.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateRoom)
		room.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteRoom)
		room.GET("/get-all", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllRoom)
	}
	course := api.Group("/course")
	{
		course.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateCourse)
		course.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateCourse)
		course.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteCourse)
		course.GET("/get-all", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllCourse)
		course.GET("/get-by-id/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.GetCourseById)
	}
	group := api.Group("/group")
	{
		group.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateGroup)
		course.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateGroup)
		course.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteGroup)
		course.GET("/get-all", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllGroup)
		course.GET("/get-by-id/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.GetGroupById)
	}
}
