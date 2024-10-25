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
// @Summary CEO
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body pb.CreateUserRequest true "User data"
// @Success 200 {object} utils.AbsResponse "Successfully created user"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
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
// @Security Bearer
// @Router /api/user/get-teachers/{isDeleted} [get]
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

// GetUserById godoc
// @Summary CEO , ADMIN , TEACHER
// @Description Retrieve user details by their unique user ID.
// @Tags user
// @Param userId path string true "User ID"
// @Produce json
// @Success 200 {object} pb.GetUserByIdResponse
// @Failure 400 {object} utils.AbsResponse
// @Router /api/user/get-user/{userId} [get]
func GetUserById(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	userId := ctx.Param("userId")
	response, err := userClient.GetUserById(ctxR, userId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// UpdateUserById godoc
// @Summary ADMIN , CEO , TEACHER
// @Description Update user details using their ID
// @Tags user
// @Accept json
// @Produce json
// @Param user body pb.UpdateUserRequest true "User Details"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Router /api/user/update [patch]
func UpdateUserById(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req := pb.UpdateUserRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := userClient.UpdateUserById(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteUserById godoc
// @Summary ADMIN , CEO
// @Description Delete a user from the system using their ID
// @Tags user
// @Param userId path string true "User ID"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Router /api/user/delete/{userId} [delete]
func DeleteUserById(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	userId := ctx.Param("userId")
	resp, err := userClient.DeleteUserById(ctxR, userId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return

}

// GetAllEmployee godoc
// @Summary ADMIN , CEO
// @Description Retrieves a list of all employees based on archive status. Restricted to ADMIN and CEO roles.
// @Tags user
// @Accept json
// @Produce json
// @Param isArchived path bool true "Filter by archive status (true=archived, false=active)"
// @Success 200 {object} pb.GetAllEmployeeResponse
// @Failure 400 {object} utils.AbsResponse "Invalid input or processing error"
// @Failure 401 {object} utils.AbsResponse "Unauthorized access"
// @Failure 408 {object} utils.AbsResponse "Request timeout"
// @Security BearerAuth
// @Router /api/user/get-all-employee/{isArchived} [get]
func GetAllEmployee(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	isArchived, err := strconv.ParseBool(ctx.Param("isArchived"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := userClient.GetAllEmployee(ctxR, isArchived)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
