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
		lead.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateLead)
		lead.GET("/get-lead-common", middleware.AuthMiddleware([]string{}, userClient), handlers.GetLeadCommon)
		lead.PUT("/update/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateLead)
		lead.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteLead)
	}
	expectation := api.Group("/expectation")
	{
		expectation.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateExpectation)
		expectation.PUT("/update/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateExpectation)
		expectation.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteExpectation)
	}
	set := api.Group("/set")
	{
		set.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateSet)
		set.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateSet)
		set.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteSet)
	}
	leadData := api.Group("/leadData")
	{
		leadData.POST("/create", middleware.AuthMiddleware([]string{}, userClient), handlers.CreateLeadData)
		leadData.PUT("/update", middleware.AuthMiddleware([]string{}, userClient), handlers.UpdateLeadData)
		leadData.DELETE("/delete/:id", middleware.AuthMiddleware([]string{}, userClient), handlers.DeleteLeadData)
	}

}
