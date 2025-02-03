package routes

import (
	client "api-gateway/internal/clients"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(r *gin.Engine, userClient *client.UserClient) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := r.Group("/api")
	{
		LeadRoutes(api, userClient)
		EducationRoutes(api, userClient)
		UserRoutes(api, userClient)
		FinanceRoutes(api, userClient)
	}
}
