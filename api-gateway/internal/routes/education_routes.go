package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func EducationRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	room := api.Group("/room")
	{
		room.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateRoom)
		room.PUT("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateRoom)
		room.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteRoom)
		room.GET("/get-all", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllRoom)
	}
	course := api.Group("/course")
	{
		course.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateCourse)
		course.PUT("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateCourse)
		course.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteCourse)
		course.GET("/get-all", handlers.GetAllCourse)
		course.GET("/get-by-id/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetCourseById)
	}
	group := api.Group("/group")
	{
		group.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateGroup)
		group.PUT("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateGroup)
		group.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteGroup)
		group.GET("/get-all/:isArchived", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllGroup)
		group.GET("/get-by-id/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.GetGroupById)
		group.GET("/get-by-course/:courseId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetGroupByCourseId)
		group.POST("/transfer-date", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.TransferLessonDate)
		group.GET("/get-by-teacher/:teacherId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetInformationByTeacher)
	}
	attendance := api.Group("/attendance")
	{
		attendance.POST("/set", middleware.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.SetAttendance)
		attendance.POST("/get-attendance", middleware.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.GetAttendance)
	}
	student := api.Group("/student")
	{
		student.GET("/get-all/:condition", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllStudent)
		student.GET("/get-student-by-id/:studentId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetStudentById)
		student.GET("/search-student/:value", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.SearchStudent)
		student.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateStudent)
		student.PUT("/update", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateStudent)
		student.DELETE("/delete/:id", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteStudent)
		student.POST("/add-to-group", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.AddStudentToGroup)
		student.PUT("/change-condition", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.ChangeConditionStudent)
		studentNote := student.Group("/note")
		{
			studentNote.GET("/get-notes/:studentId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetNotesByStudent)
			studentNote.POST("/create", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateNoteForStudent)
			studentNote.DELETE("/delete/:noteId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteStudentNote)
		}
	}
	history := api.Group("/history")
	{
		history.GET("/group/:groupId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetHistoryGroup)
		history.GET("/student/:studentId", middleware.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetHistoryStudent)
	}
}
