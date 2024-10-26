package middleware

import (
	client "api-gateway/internal/clients"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRoles []string, userClient *client.UserClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//authHeader := ctx.GetHeader("Authorization")
		//if authHeader == "" {
		//	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		//	ctx.Abort()
		//	return
		//}
		//
		//const bearerPrefix = "Bearer "
		//if !strings.HasPrefix(authHeader, bearerPrefix) {
		//	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		//	ctx.Abort()
		//	return
		//}
		//token := strings.TrimPrefix(authHeader, bearerPrefix)
		//
		//user, err := userClient.ValidateToken(token, requiredRoles)
		//if err != nil {
		//	ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid or insufficient permissions required role => %v", requiredRoles)})
		//	ctx.Abort()
		//	return
		//}
		//
		//ctx.Set("user", user)
		//ctx.Next()
	}
}
