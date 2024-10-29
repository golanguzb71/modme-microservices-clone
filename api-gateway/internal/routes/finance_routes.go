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
		category := finance.Group("/category")
		{
			category.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateCategory)
			category.DELETE("/delete/:categoryId", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteCategory)
			category.GET("/get-all", middleware.AuthMiddleware([]string{}, userClient), handlers.GetAllCategories)
		}
	}

}
