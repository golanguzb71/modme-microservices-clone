package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/etc"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func FinanceRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	finance := api.Group("/finance")
	{
		discount := finance.Group("/discount")
		{
			discount.GET("/get-all-by-group/:groupId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllDiscountInformationByGroup)
			discount.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateDiscount)
			discount.DELETE("/delete", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteDiscount)
			discount.GET("/history/:userId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetHistoryDiscount)
		}
		category := finance.Group("/category")
		{
			category.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateCategory)
			category.DELETE("/delete/:categoryId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteCategory)
			category.GET("/get-all", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllCategories)
		}
		expense := finance.Group("/expense")
		{
			expense.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateExpense)
			expense.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteExpense)
			expense.GET("/get-all-information/:from/:to", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllInformation)
			expense.GET("/get-chart-diagram/:from/:to", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetChartDiagram)
		}
		payment := finance.Group("/payment")
		{
			payment.POST("/student/add", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.PaymentAdd)
			payment.POST("/student/return", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.PaymentReturn)
			payment.PATCH("/student/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.PaymentUpdate)
			payment.GET("/student/get-monthly-status/:studentId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetMonthlyStatusPayment)
			payment.GET("/get-all-payments/:studentId/:month", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllPayments)
			payment.GET("/payment-take-off/:from/:to", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllTakeOffPayment)
			payment.GET("/payment-take-off/chart/:from/:to", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetPaymentTakeOffChart)
			payment.POST("/all-student-payments", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllStudentPayment)
			payment.POST("/all-student-payments/chart", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllPaymentsStudentChart)
			payment.GET("/get-all-debts/:page/:size", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllDebtsInformation)
		}
		salary := finance.Group("/salary")
		{
			salary.GET("/teacher-all", etc.AuthMiddleware([]string{"CEO"}, userClient), handlers.GetSalaryAllTeacher)
			salary.POST("/teacher-add", etc.AuthMiddleware([]string{"CEO"}, userClient), handlers.AddSalaryTeacher)
			salary.DELETE("/delete/:teacherID", etc.AuthMiddleware([]string{"CEO"}, userClient), handlers.DeleteTeacherSalary)
			salary.GET("/calculate/:from/:to", etc.AuthMiddleware([]string{"CEO"}, userClient), handlers.CalculateSalary)
		}
	}
}
