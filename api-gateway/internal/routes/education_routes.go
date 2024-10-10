package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func EducationRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	room := api.Group("/room")
	{
		room.POST("/create", handlers.CreateRoom)
		room.PUT("/update", handlers.UpdateRoom)
		room.DELETE("/delete/:id", handlers.DeleteRoom)
		room.GET("/get-all", handlers.GetAllRoom)
	}
}
