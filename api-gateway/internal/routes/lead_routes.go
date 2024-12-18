package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LeadRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	lead := api.Group("/lead")
	{
		lead.POST("/create", handlers.CreateLead)
		lead.POST("/get-lead-common", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetLeadCommon)
		lead.PUT("/update/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateLead)
		lead.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteLead)
		lead.GET("/get-all", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllLead)
		lead.GET("/get-lead-reports", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetLeadReports)
	}
	expectation := api.Group("/expectation")
	{
		expectation.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateExpectation)
		expectation.PUT("/update/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateExpectation)
		expectation.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteExpectation)
	}
	set := api.Group("/set")
	{
		set.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateSet)
		set.PUT("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateSet)
		set.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteSet)
		set.PATCH("/change-to-group", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.ChangeToSet)
	}
	leadData := api.Group("/leadData")
	{
		leadData.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateLeadData)
		leadData.PUT("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateLeadData)
		leadData.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteLeadData)
		leadData.PATCH("/change-lead-data", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.ChangeLeadData)
	}
}
