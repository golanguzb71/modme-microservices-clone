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

// CreateRoom godoc
// @Summary ADMIN , CEO
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
// @Summary ADMIN , CEO
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
// @Summary ADMIN , CEO
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
// @Summary ADMIN , CEO
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
// @Summary ADMIN , CEO
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
// @Summary ADMIN , CEO
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
// @Summary ADMIN , CEO
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
// @Summary ALL
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
// @Summary ADMIN , CEO
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

// CreateGroup godoc
// @Summary ADMIN , CEO
// @Description Create a new group with provided details.
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group body pb.CreateGroupRequest true "Group Data"
// @Success 200 {object} utils.AbsResponse "Group successfully created"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/create [post]
func CreateGroup(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req pb.CreateGroupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.CreateGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateGroup godoc
// @Summary ADMIN , CEO
// @Description Update details of an existing group.
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group body pb.GetUpdateGroupAbs true "Group Data"
// @Success 200 {object} utils.AbsResponse "Group successfully updated"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/update [put]
func UpdateGroup(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req *pb.GetUpdateGroupAbs
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.UpdateGroup(ctxR, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteGroup godoc
// @Summary ADMIN , CEO
// @Description Delete a group by its ID.
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} utils.AbsResponse "Group successfully deleted"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/delete/{id} [delete]
func DeleteGroup(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.DeleteGroup(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllGroup godoc
// @Summary ADMIN , CEO
// @Description Retrieve a list of all groups.
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param isArchived path bool true "Is Archived" example(true)
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Success 200 {array} pb.GetGroupsResponse "List of groups"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/get-all/{isArchived} [get]
func GetAllGroup(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	isArchived := ctx.Param("isArchived")
	parseBool, err := strconv.ParseBool(isArchived)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid isArchived parameter"})
		return
	}
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 32)
	if err != nil || size < 1 {
		size = 10
	}
	resp, err := educationClient.GetAllGroup(ctxR, parseBool, int32(page), int32(size))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
}

// GetGroupById godoc
// @Summary ADMIN
// @Description Retrieve details of a group by its ID.
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} pb.GetGroupAbsResponse "Group details"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/get-by-id/{id} [get]
func GetGroupById(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.GetGroupById(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
	return
}

// SetAttendance godoc
// @Summary TEACHER
// @Description Record attendance for a student in a group on a specific date.
// @Tags attendance
// @Produce json
// @Security BearerAuth
// @Param attendance body pb.SetAttendanceRequest true "Attendance details"
// @Success 200 {object} utils.AbsResponse "Attendance recorded successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/attendance/set [post]
func SetAttendance(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req pb.SetAttendanceRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.SetAttendanceByGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAttendance godoc
// @Summary ADMIN , TEACHER
// @Description Retrieve attendance records for students in a group over a specified date range.
// @Tags attendance
// @Produce json
// @Security BearerAuth
// @Param attendance body pb.GetAttendanceRequest true "Group ID and date range"
// @Success 200 {object} pb.GetAttendanceResponse "Attendance records"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/attendance/get-attendance [post]
func GetAttendance(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req pb.GetAttendanceRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.GetAttendanceByGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
	return
}

// GetGroupByCourseId godoc
// @Summary ADMIN, TEACHER
// @Description Retrieve groups associated with a specific course ID.
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param courseId path string true "Course ID"
// @Success 200 {object} pb.GetGroupsByCourseResponse "Group details"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/get-by-course/{courseId} [get]
func GetGroupByCourseId(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	courseId := ctx.Param("courseId")
	resp, err := educationClient.GetGroupByCourseId(ctxR, courseId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
	return
}

// GetAllStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security BearerAuth
// @Param condition path string true "Condition" Enums(ARCHIVED, ACTIVE)
// @Param page query string false "Page number"
// @Param size query string false "Page size"
// @Success 200 {object} pb.GetAllStudentResponse "List of students"
// @Failure 400 {object} utils.AbsResponse "Invalid condition"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/get-all/{condition} [get]
func GetAllStudent(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	condition := ctx.Param("condition")
	if condition != "ARCHIVED" && condition != "ACTIVE" {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid condition")
		return
	}
	page := ctx.Query("page")
	size := ctx.Query("size")
	response, err := educationClient.GetAllStudent(ctxR, condition, page, size)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// CreateStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security BearerAuth
// @Param request body pb.CreateStudentRequest true "Student details"
// @Success 200 {object} utils.AbsResponse "Created student details"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/create [post]
func CreateStudent(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req := pb.CreateStudentRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.CreatedBy = "c1d6503f-31dc-4f99-b61f-2e4ebc7a7639"
	response, err := educationClient.CreateStudent(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// AddStudentToGroup godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security BearerAuth
// @Param request body pb.AddToGroupRequest true "Add student to group details"
// @Success 200 {object} utils.AbsResponse "Success message"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/add-to-group [post]
func AddStudentToGroup(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req := pb.AddToGroupRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.CreatedBy = "c1d6503f-31dc-4f99-b61f-2e4ebc7a7639"
	response, err := educationClient.AddStudentToGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, response.Status, response.Message)
	return
}

// UpdateStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security BearerAuth
// @Param request body pb.UpdateStudentRequest true "Updated student details"
// @Success 200 {object} utils.AbsResponse "Update success message"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/update [put]
func UpdateStudent(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req := pb.UpdateStudentRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response, err := educationClient.UpdateStudent(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, response.Status, response.Message)
	return
}

// DeleteStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} utils.AbsResponse "Delete success message"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/delete/{id} [delete]
func DeleteStudent(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := ctx.Param("id")
	response, err := educationClient.DeleteStudent(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, response.Status, response.Message)
	return
}

// GetStudentById godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Success 200 {object} pb.GetAllStudentResponse "List of students"
// @Failure 400 {object} utils.AbsResponse "Invalid condition"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/get-student-by-id/{studentId} [get]
func GetStudentById(ctx *gin.Context) {

}

func GetNotesByStudent(ctx *gin.Context) {

}

func CreateNoteForStudent(ctx *gin.Context) {

}

func DeleteStudentNote(ctx *gin.Context) {

}
