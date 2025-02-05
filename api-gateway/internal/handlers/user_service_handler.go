package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/etc"
	"api-gateway/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
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
	ctxR, cancel := etc.NewTimoutContext(ctx)
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
	isDeleted, err := strconv.ParseBool(ctx.Param("isDeleted"))
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
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
	ctxR, cancel := etc.NewTimoutContext(ctx)
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
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.UpdateUserRequest{}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	if user.Role == "TEACHER" || user.Role == "EMPLOYEE" {
		req.Id = user.Id
	}
	err = ctx.ShouldBindJSON(&req)
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
// @Security Bearer
// @Router /api/user/delete/{userId} [delete]
func DeleteUserById(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
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
// @Security Bearer
// @Router /api/user/get-all-employee/{isArchived} [get]
func GetAllEmployee(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
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

// Login godoc
// @Summary ALL
// @Description Authenticate a user and return a token upon successful login.
// @Tags user
// @Accept json
// @Produce json
// @Param LoginRequest body pb.LoginRequest true "Login credentials"
// @Success 200 {object} pb.LoginResponse "Successful login"
// @Failure 400 {object} utils.AbsResponse "Bad request - Invalid JSON or login failure"
// @Router /api/user/login [post]
func Login(ctx *gin.Context) {
	ctxR := context.Background()
	req := pb.LoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(&req)
	resp, err := userClient.Login(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetMyInformation godoc
// @Summary ADMIN , CEO , TEACHER
// @Description Retrieve the information of the authenticated user. Admin and CEO roles receive basic user info, while Teachers get additional group information.
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User information with optional groups for Teachers"
// @Failure 401 {object} utils.AbsResponse "Unauthorized - Invalid or missing authentication"
// @Failure 500 {object} utils.AbsResponse "Internal server error while retrieving teacher's group information"
// @Security Bearer
// @Router /api/user/get-my-profile [get]
func GetMyInformation(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	if user.Role == "ADMIN" || user.Role == "CEO" || user.Role == "SUPER_CEO" {
		ctx.JSON(http.StatusOK, user)
		return
	}
	if user.Role == "TEACHER" {
		groups, err := educationClient.GetInformationByTeacher(ctxR, user.Id, false)
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		response := map[string]interface{}{
			"user_info": user,
			"groups":    groups,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	utils.RespondError(ctx, http.StatusUnauthorized, "aborted: required role ADMIN, CEO, SUPER_CEO or TEACHER")
}

// GetAllStaff godoc
// @Summary ADMIN , CEO
// @Description Retrieve a list of all staff members, filtered by their archived status.
// @Tags user
// @Accept json
// @Produce json
// @Param isArchived path boolean true "Boolean to filter archived or active staff; true for archived, false for active"
// @Success 200 {array} pb.GetAllStuffResponse "List of staff members based on archived status"
// @Failure 400 {object} utils.AbsResponse "Bad request - Invalid 'isArchived' parameter"
// @Failure 500 {object} utils.AbsResponse "Internal server error during data retrieval"
// @Security Bearer
// @Router /api/user/get-all-staff/{isArchived} [get]
func GetAllStaff(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	isArchived, err := strconv.ParseBool(ctx.Param("isArchived"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := userClient.GetAllStuff(ctxR, isArchived)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetUserHistoryById retrieves the history for a specific user by user ID.
// @Summary ADMIN, CEO, TEACHER
// @Description Fetches the history data for a given user ID.
// @Tags user
// @Param userId path string true "User ID"
// @Success 200 {object} pb.GetHistoryByUserIdResponse
// @Failure 400 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/user/history/{userId} [get]
func GetUserHistoryById(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	userId := ctx.Param("userId")
	resp, err := userClient.GetHistoryByUserId(ctxR, userId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// UpdateUserPassword updates a user's password.
// @Summary CEO
// @Description Updates the password of a user specified by the userId.
// @Tags user
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param password path string true "New Password"
// @Success 200 {object} utils.AbsResponse "Password updated successfully"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Security Bearer
// @Router /api/user/update-password/{userId}/{password} [put]
func UpdateUserPassword(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	userId := ctx.Param("userId")
	password := ctx.Param("password")
	resp, err := userClient.UpdateUserPassword(ctxR, userId, password)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// CreateUserForCompany godoc
// @Summary SUPER_CEO
// @Description Create a new user associated with a company
// @Tags company-user
// @Accept json
// @Produce json
// @Param companyId query string true "Company ID"
// @Param user body pb.CreateUserRequest true "User data"
// @Success 200 {object} utils.AbsResponse "Successfully created user"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/company-user/create [post]
func CreateUserForCompany(ctx *gin.Context) {
	companyId := ctx.Query("companyId")
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctxR := metadata.NewOutgoingContext(context.Background(), md)

	var userForCompany pb.CreateUserRequest
	if err := ctx.ShouldBindJSON(&userForCompany); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := userClient.CreateUser(ctxR, &userForCompany)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetUserByIdForCompany godoc
// @Summary SUPER_CEO
// @Description Retrieve details of a user by their ID and company ID
// @Tags company-user
// @Accept json
// @Produce json
// @Param companyId query string true "Company ID"
// @Param userId path string true "User ID"
// @Success 200 {object} pb.GetUserByIdResponse "Successfully retrieved user"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 404 {object} utils.AbsResponse "User not found"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/company-user/get-user/{userId} [get]
func GetUserByIdForCompany(ctx *gin.Context) {
	companyId := ctx.Query("companyId")
	userId := ctx.Param("userId")
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctxR := metadata.NewOutgoingContext(ctx, md)
	response, err := userClient.GetUserById(ctxR, userId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// UpdateUserbyIdForCompany godoc
// @Summary SUPER_CEO
// @Description Update user details by their ID and company ID
// @Tags company-user
// @Accept json
// @Produce json
// @Param companyId query string true "Company ID"
// @Param user body pb.UpdateUserRequest true "Updated user data"
// @Success 200 {object} utils.AbsResponse "Successfully updated user"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 404 {object} utils.AbsResponse "User not found"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/company-user/update [patch]
func UpdateUserbyIdForCompany(ctx *gin.Context) {
	req := pb.UpdateUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	companyId := ctx.Query("companyId")
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctxR := metadata.NewOutgoingContext(ctx, md)
	resp, err := userClient.UpdateUserById(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteUserByIdForCompany godoc
// @Summary SUPER_CEO
// @Description Delete a user by their ID and associated company ID
// @Tags company-user
// @Accept json
// @Produce json
// @Param companyId query string true "Company ID"
// @Param userId path string true "User ID"
// @Success 200 {object} utils.AbsResponse "Successfully deleted user"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 404 {object} utils.AbsResponse "User not found"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/company-user/delete/{userId} [delete]
func DeleteUserByIdForCompany(ctx *gin.Context) {
	companyId := ctx.Query("companyId")
	userId := ctx.Param("userId")
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctxR := metadata.NewOutgoingContext(ctx, md)
	resp, err := userClient.DeleteUserById(ctxR, userId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}
