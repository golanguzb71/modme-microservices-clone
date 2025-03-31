package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/etc"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	user := api.Group("/user")
	{
		user.POST("/login", handlers.Login)
		user.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateUser)
		user.GET("/get-teachers/:isDeleted", etc.AuthMiddleware([]string{"ADMIN", "CEO", "FINANCIST"}, userClient), handlers.GetTeachers)
		user.GET("/get-user/:userId", etc.AuthMiddleware([]string{"ADMIN", "CEO", "FINANCIST"}, userClient), handlers.GetUserById)
		user.PATCH("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO", "FINANCIST"}, userClient), handlers.UpdateUserById)
		user.DELETE("/delete/:userId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteUserById)
		user.GET("/get-all-employee/:isArchived", etc.AuthMiddleware([]string{"ADMIN", "CEO", "FINANCIST"}, userClient), handlers.GetAllEmployee)
		user.GET("/get-my-profile", etc.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER", "SUPER_CEO", "FINANCIST"}, userClient), handlers.GetMyInformation)
		user.GET("/get-all-staff/:isArchived", etc.AuthMiddleware([]string{"ADMIN", "CEO", "FINANCIST"}, userClient), handlers.GetAllStaff)
		user.GET("/history/:userId", etc.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER", "FINANCIST"}, userClient), handlers.GetUserHistoryById)
		user.PUT("/update-password/:userId/:password", etc.AuthMiddleware([]string{"CEO", "FINANCIST"}, userClient), handlers.UpdateUserPassword)
	}

	companyUser := api.Group("/company-user")
	{
		companyUser.POST("/create", handlers.CreateUserForCompany)
		companyUser.GET("/get-user/:userId", handlers.GetUserByIdForCompany)
		companyUser.PATCH("/update", handlers.UpdateUserbyIdForCompany)
		companyUser.DELETE("/delete/:userId", handlers.DeleteUserByIdForCompany)
	}
}
