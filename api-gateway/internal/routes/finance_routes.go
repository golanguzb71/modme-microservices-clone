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
		discount.GET("/get-all-by-group/:groupId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllDiscountInformationByGroup)
		discount.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateDiscount)
		discount.DELETE("/delete", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteDiscount)
		category := finance.Group("/category")
		{
			category.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateCategory)
			category.DELETE("/delete/:categoryId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteCategory)
			category.GET("/get-all", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllCategories)
		}
		expense := finance.Group("/expense")
		{
			expense.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateExpense)
			expense.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteExpense)
			expense.GET("/get-all-information/:from/:to", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllInformation)
			expense.GET("/get-chart-diagram", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetChartDiagram)
		}

	}
}
