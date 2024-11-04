package utils

import (
	"api-gateway/grpc/proto/pb"
	"errors"
	"github.com/gin-gonic/gin"
)

// AbsResponse represents an error API response
type AbsResponse struct {
	Status  int32  `json:"statusCode"`
	Message string `json:"message"`
}

func RespondSuccess(ctx *gin.Context, statusCode int32, message string) {
	ctx.JSON(int(statusCode), AbsResponse{
		Status:  statusCode,
		Message: message,
	})
}

func RespondError(ctx *gin.Context, statusCode int32, message string) {
	ctx.JSON(int(statusCode), AbsResponse{
		Status:  statusCode,
		Message: message,
	})
}

func GetUserFromContext(c *gin.Context) (*pb.GetUserByIdResponse, error) {
	var userInterface any
	var exists bool
	userInterface, exists = c.Get("user")
	if !exists {
		return nil, errors.New("user not found in Gin context")
	}

	user, ok := userInterface.(*pb.GetUserByIdResponse)
	if !ok {
		return nil, errors.New("error while converting user")
	}
	return user, nil
}
