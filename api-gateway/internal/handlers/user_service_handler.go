package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreateUser godoc
// @Summary Create a new user
// @Description CEO
// @Tags user
// @Accept json
// @Produce json
// @Param user body pb.CreateUserRequest true "User data"
// @Success 200 {object} utils.AbsResponse "Successfully created user"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/user/create [post]
func CreateUser(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req *pb.CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := userClient.CreateUser(ctxR, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetTeachers godoc
// @Summary ADMIN,CEO
// @Description Fetches a list of teachers based on the deletion status
// @Tags user
// @Accept json
// @Produce json
// @Param isDeleted path bool true "Deletion status (true/false)"
// @Success 200 {object} pb.GetTeachersResponse "List of teachers"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/user/teachers/{isDeleted} [get]
func GetTeachers(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	isDeleted, err := strconv.ParseBool(ctx.Param("isDeleted"))
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	resp, err := userClient.GetTeachers(ctxR, isDeleted)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
