package routes

import (
	client "api-gateway/internal/clients"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(r *gin.Engine, userClient *client.UserClient) {
	r.GET("/swag", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		LeadRoutes(api, userClient)
	}
}
