package main

import (
	"api-gateway/config"
	_ "api-gateway/docs"
	"api-gateway/grpc"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// @title Modme Swagger

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	router := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	grpcClients := grpc.InitializeGrpcClients(cfg)
	handlers.InitClients(grpcClients)
	routes.SetUpRoutes(router, grpcClients.UserClient)

	port := cfg.Server.Port
	log.Printf("Starting api gateway on port %s", port)
	if err = router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
