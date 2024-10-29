package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	user := api.Group("/user")
	{
		user.POST("/login", handlers.Login)
		user.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateUser)
		user.GET("/get-teachers/:isDeleted", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetTeachers)
		user.GET("/get-user/:userId", handlers.GetUserById)
		user.PATCH("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateUserById)
		user.DELETE("/delete/:userId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteUserById)
		user.GET("/get-all-employee/:isArchived", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllEmployee)
		user.GET("/get-my-profile", middleware.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.GetMyInformation)
	}
}
