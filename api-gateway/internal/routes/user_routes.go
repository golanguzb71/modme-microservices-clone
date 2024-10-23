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
		user.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateUser)
		user.GET("/get-teachers/:isDeleted", middleware.AuthMiddleware([]string{}, userClient), handlers.GetTeachers)
	}
}
