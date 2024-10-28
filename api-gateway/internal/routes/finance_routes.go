package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func FinanceRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	finance := api.Group("/finance")
	{
		discount := finance.Group("/discount")
		discount.GET("/get-all-by-group/:groupId", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllDiscountInformationByGroup)
		discount.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateDiscount)
		discount.DELETE("/delete", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteDiscount)

	}
}
