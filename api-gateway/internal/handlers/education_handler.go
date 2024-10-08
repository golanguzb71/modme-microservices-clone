package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CreateRoom godoc
// @Summary ADMIN
// @Description Create a new room based on the provided request data
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.CreateRoomRequest true "Request to create a room"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/create [post]
func CreateRoom(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req pb.CreateRoomRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	resp, err := educationClient.CreateRoom(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateRoom godoc
// @Summary ADMIN
// @Description Update the details of an existing room based on the provided request data
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request query pb.AbsRoom true "Request to update room"
// @Success 200 {object} utils.AbsResponse
// @Failure 422 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/update [put]
func UpdateRoom(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req pb.AbsRoom
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	resp, err := educationClient.UpdateRoom(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteRoom godoc
// @Summary ADMIN
// @Description Delete a room by its ID
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 200 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/delete/{id} [delete]
func DeleteRoom(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.DeleteRoom(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllRoom godoc
// @Summary ADMIN
// @Description Retrieve all rooms
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pb.GetUpdateRoomAbs
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/get-all [get]
func GetAllRoom(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rooms, err := educationClient.GetRoom(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rooms)
	return
}

// CreateCourse godoc
// @Summary ADMIN
// @Description Create a new course based on the provided request data
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.CreateCourseRequest true "Request to create a course"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/create [post]
func CreateCourse(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req pb.CreateCourseRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.CreateCourse(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateCourse godoc
// @Summary ADMIN
// @Description Update the details of an existing course based on the provided request data
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.AbsCourse true "Request to update course"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/update [put]
func UpdateCourse(ctx *gin.Context) {
	var req pb.AbsCourse
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.UpdateCourse(ctx, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteCourse godoc
// @Summary ADMIN
// @Description Delete a course by its ID
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Success 200 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/delete/{id} [delete]
func DeleteCourse(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.DeleteCourse(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllCourse godoc
// @Summary ADMIN
// @Description Retrieve all courses
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pb.GetUpdateCourseAbs
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/get-all [get]
func GetAllCourse(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rooms, err := educationClient.GetCourse(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rooms)
	return
}

// GetCourseById godoc
// @Summary Retrieve course by ID (ADMIN)
// @Description Retrieves a course by its ID for admin users.
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Success 200 {object} pb.AbsCourse "Successful response with course details"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/course/get-by-id/{id} [get]
func GetCourseById(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.GetCourseById(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
