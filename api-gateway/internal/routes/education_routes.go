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
		group.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateGroup)
		group.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteGroup)
		group.GET("/get-all/:isArchived", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllGroup)
		group.GET("/get-by-id/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.GetGroupById)
		group.GET("/get-by-course/:courseId", middleware.AuthMiddleware([]string{}, userClient), handlers.GetGroupByCourseId)
	}
	attendance := api.Group("/attendance")
	{
		attendance.POST("/set", middleware.AuthMiddleware([]string{}, userClient), handlers.SetAttendance)
		attendance.POST("/get-attendance", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAttendance)
	}
	student := api.Group("/student")
	{
		student.GET("/get-all/:condition", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllStudent)
		student.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateStudent)
		student.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateStudent)
		student.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteStudent)
		student.POST("/add-to-group", middleware.AuthMiddleware([]string{}, userClient), handlers.AddStudentToGroup)
	}
}
