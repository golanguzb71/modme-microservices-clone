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
		discount.GET("/history/:userId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetHistoryDiscount)
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
			expense.GET("/get-chart-diagram/:from/:to", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetChartDiagram)
		}
		payment := finance.Group("/payment")
		{
			payment.POST("/student/add", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.PaymentAdd)
			payment.POST("/student/return", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.PaymentReturn)
			payment.PATCH("/student/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.PaymentUpdate)
			payment.GET("/student/get-monthly-status/:studentId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetMonthlyStatusPayment)
			payment.GET("/get-all-payments/:studentId/:month", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllPayments)
		}
	}
}
