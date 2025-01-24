package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/etc"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func EducationRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	api.GET("/common-information-company", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetCommonInformationCompany)
	api.GET("/get-chart-income", etc.AuthMiddleware([]string{"CEO"}, userClient), handlers.GetChartIncome)
	api.GET("/get-table-groups", etc.AuthMiddleware([]string{"CEO", "ADMIN", "TEACHER"}, userClient), handlers.GetTableGroups)
	image := api.Group("/image")
	{
		image.POST("/upload", handlers.UploadImage)
		image.GET("/get-image", handlers.GetImage)
	}
	company := api.Group("/company")
	{
		company.GET("/subdomain/:domain", handlers.GetCompanyBySubdomain)
		company.POST("/create", handlers.CompanyCreate)
		company.GET("/get-all", handlers.GetAllCompanies)
		company.PUT("/update", handlers.CompanyUpdate)
		tariff := company.Group("/tariff")
		tariff.POST("/create", handlers.TariffCreate)
		tariff.GET("/get-all", handlers.TariffGetAll)
		tariff.PUT("/update", handlers.TariffUpdate)
		tariff.DELETE("/delete/:id", handlers.TariffDelete)
		finance := company.Group("/finance")
		finance.POST("/create", handlers.FinanceCreate)
		finance.DELETE("/delete", handlers.FinanceDelete)
		finance.POST("/get-all", handlers.FinanceGetAll)
		finance.POST("/get-by-company", handlers.FinanceGetByCompany)
	}
	room := api.Group("/room")
	{
		room.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateRoom)
		room.PUT("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateRoom)
		room.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteRoom)
		room.GET("/get-all", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllRoom)
	}
	course := api.Group("/course")
	{
		course.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateCourse)
		course.PUT("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateCourse)
		course.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteCourse)
		course.GET("/get-all", etc.AuthMiddleware([]string{"ADMIN", "CEO", "SUPER_CEO"}, userClient), handlers.GetAllCourse)
		course.GET("/get-by-id/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetCourseById)
	}
	group := api.Group("/group")
	{
		group.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateGroup)
		group.PUT("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateGroup)
		group.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteGroup)
		group.GET("/get-all/:isArchived", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllGroup)
		group.GET("/get-by-id/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.GetGroupById)
		group.GET("/get-by-course/:courseId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetGroupByCourseId)
		group.POST("/transfer-date", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.TransferLessonDate)
		group.GET("/get-by-teacher/:teacherId", etc.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.GetInformationByTeacher)
		group.GET("/left-after-trial/:from/:to", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.LeftAfterTrial)
	}
	attendance := api.Group("/attendance")
	{
		attendance.POST("/set", etc.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.SetAttendance)
		attendance.POST("/get-attendance", etc.AuthMiddleware([]string{"ADMIN", "CEO", "TEACHER"}, userClient), handlers.GetAttendance)
	}
	student := api.Group("/student")
	{
		student.GET("/get-all/:condition", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetAllStudent)
		student.GET("/get-student-by-id/:studentId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetStudentById)
		student.GET("/search-student/:value", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.SearchStudent)
		student.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateStudent)
		student.PUT("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateStudent)
		student.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteStudent)
		student.POST("/add-to-group", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.AddStudentToGroup)
		student.PUT("/change-condition", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.ChangeConditionStudent)
		studentNote := student.Group("/note")
		{
			studentNote.GET("/get-notes/:studentId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetNotesByStudent)
			studentNote.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateNoteForStudent)
			studentNote.DELETE("/delete/:noteId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteStudentNote)
		}
	}
	history := api.Group("/history")
	{
		history.GET("/group/:groupId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetHistoryGroup)
		history.GET("/student/:studentId", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetHistoryStudent)
	}

}
