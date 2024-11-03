package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"time"
)

// GetAllDiscountInformationByGroup godoc
// @Summary ADMIN , CEO
// @Description Retrieves all discount information for a specific group ID
// @Tags discount
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 200 {object} pb.GetInformationDiscountResponse "Success response"
// @Failure 500 {object} utils.AbsResponse "Internal Server Error"
// @Security Bearer
// @Router /api/finance/discount/get-all-by-group/{groupId} [get]
func GetAllDiscountInformationByGroup(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	groupId := ctx.Param("groupId")
	resp, err := financeClient.GetDiscountsInformationByGroupId(ctxR, groupId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// CreateDiscount godoc
// @Summary CEO , ADMIN
// @Description Creates a new discount for a specified group and student
// @Tags discount
// @Accept json
// @Produce json
// @Param request body pb.AbsDiscountRequest true "Discount Request"
// @Success 200 {object} utils.AbsResponse "Success response with status and message"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Failure 500 {object} utils.AbsResponse "Internal Server Error"
// @Security Bearer
// @Router /api/finance/discount/create [post]
func CreateDiscount(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb.AbsDiscountRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := financeClient.CreateDiscount(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteDiscount godoc
// @Summary ADMIN , CEO
// @Description Deletes a discount for a specific group and student
// @Tags discount
// @Accept json
// @Produce json
// @Param groupId query string true "Group ID"
// @Param studentId query string true "Student ID"
// @Security Bearer
// @Success 200 {object} utils.AbsResponse "Success response with status and message"
// @Failure 500 {object} utils.AbsResponse "Internal Server Error"
// @Router /api/finance/discount/delete [delete]
func DeleteDiscount(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	groupId := ctx.Query("groupId")
	studentId := ctx.Query("studentId")
	resp, err := financeClient.DeleteDiscount(ctxR, groupId, studentId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// CreateCategory godoc
// @Summary      ADMIN , CEO
// @Description  Creates a new category with the provided name and description
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        category  body      pb.CreateCategoryRequest  true  "Category Data"
// @Success      200       {object}  utils.AbsResponse
// @Failure      400       {object}  utils.AbsResponse
// @Failure      500       {object}  utils.AbsResponse
// @Security Bearer
// @Router       /api/finance/category/create [post]
func CreateCategory(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb.CreateCategoryRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := financeClient.CreateCategory(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteCategory godoc
// @Summary      ADMIN , CEO
// @Description  Deletes a category by its ID
// @Tags         category
// @Produce      json
// @Param        categoryId  path      string  true  "Category ID"
// @Success      200         {object}  utils.AbsResponse
// @Failure      400         {object}  utils.AbsResponse
// @Failure      500         {object}  utils.AbsResponse
// @Security Bearer
// @Router       /api/finance/category/delete/{categoryId} [delete]
func DeleteCategory(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	categoryId := ctx.Param("categoryId")
	resp, err := financeClient.DeleteCategory(ctxR, categoryId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllCategories godoc
// @Summary      ADMIN , CEO
// @Description  Retrieves all categories
// @Tags         category
// @Produce      json
// @Success      200  {object}  pb.GetAllCategoryRequest
// @Failure      500  {object}  utils.AbsResponse
// @Security Bearer
// @Router       /api/finance/category/get-all [get]
func GetAllCategories(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := financeClient.GetAllCategories(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// CreateExpense godoc
// @Summary      ADMIN , CEO
// @Description  Creates a new expense entry with details provided in the request body.
// @Tags         expense
// @Accept       json
// @Produce      json
// @Param        expense  body      pb.CreateExpenseRequest true  "Expense details"
// @Success      200      {object}  utils.AbsResponse
// @Failure      500      {object}  utils.AbsResponse
// @Security Bearer
// @Router       /api/finance/expense/create [post]
func CreateExpense(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb.CreateExpenseRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	req.CreatedById = user.Id
	resp, err := financeClient.CreateExpense(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteExpense godoc
// @Summary      CEO
// @Description  Deletes an expense entry by ID.
// @Tags         expense
// @Param        id       path      string true  "Expense ID"
// @Produce      json
// @Success      200      {object}  utils.AbsResponse
// @Failure      500      {object}  utils.AbsResponse
// @Security Bearer
// @Router       /api/finance/expense/delete/{id} [delete]
func DeleteExpense(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id := ctx.Param("id")
	resp, err := financeClient.DeleteExpense(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllInformation godoc
// @Summary      ADMIN , CEO
// @Description  Retrieves expenses with optional filters, pagination, and date range.
// @Tags         expense
// @Produce      json
// @Param        from     path      string true  "Start date"
// @Param        to       path      string true  "End date"
// @Param        page     query     int    true  "Page number"
// @Param        size     query     int    true  "Page size"
// @Param        id       query     string false "ID to filter by user or creator"
// @Param        type     query     string false "Filter type (USER or CATEGORY)"
// @Success      200      {object}  pb.GetAllExpenseResponse
// @Failure      400      {object}  utils.AbsResponse  "Invalid request parameter"
// @Failure      500      {object}  utils.AbsResponse
// @Security Bearer
// @Router       /api/finance/expense/get-all-information/{from}/{to} [get]
func GetAllInformation(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	from := ctx.Param("from")
	to := ctx.Param("to")
	page, err := strconv.ParseInt(ctx.Query("page"), 10, 32)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	size, err := strconv.ParseInt(ctx.Query("size"), 10, 32)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id := ctx.Query("id")
	idType := ctx.Query("type")
	if idType == "USER" || idType == "CATEGORY" || idType == "" {
		resp, err := financeClient.GetAllInformation(ctxR, id, idType, page, size, from, to)
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, resp)
		return
	} else {
		utils.RespondError(ctx, http.StatusInternalServerError, "invalid idType")
		return
	}
}

// GetChartDiagram godoc
// @Summary ADMIN , CEO
// @Description Queries the finance service for expense data between 'from' and 'to' dates, returning it as a chart diagram.
// @Tags expense
// @Accept json
// @Produce json
// @Param from path string true "Start date for the chart (YYYY-MM-DD)"
// @Param to path string true "End date for the chart (YYYY-MM-DD)"
// @Success 200 {object} pb.GetAllExpenseDiagramResponse "Chart data successfully retrieved"
// @Failure 409 {object} utils.AbsResponse "Conflict or error retrieving chart data"
// @Security Bearer
// @Router /api/finance/expense/get-chart-diagram/{from}/{to} [get]
func GetChartDiagram(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	from := ctx.Param("from")
	to := ctx.Param("to")
	resp, err := financeClient.GetExpenseChartDiagram(from, to, ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetHistoryDiscount retrieves the discount history for a specific user.
// @Summary ADMIN , CEO
// @Description Retrieves the discount history for a user by their ID
// @Tags discount
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} pb.GetHistoryDiscountResponse "Success"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Failure 500 {object} utils.AbsResponse "Internal Server Error"
// @Security Bearer
// @Router /api/finance/discount/history/{userId} [get]
func GetHistoryDiscount(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	userId := ctx.Param("userId")
	resp, err := financeClient.GetHistoryDiscount(userId, ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// PaymentAdd godoc
// @Summary ADMIN , CEO
// @Description Add a payment for a student
// @Tags payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body pb.PaymentAddRequest true "Payment Add Request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Router /api/finance/payment/student/add [post]
func PaymentAdd(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	req := pb.PaymentAddRequest{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	amount, err := decimal.NewFromString(req.Sum)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if amount.LessThan(decimal.NewFromInt(10000)) {
		utils.RespondError(ctx, http.StatusBadRequest, "min 10 000 USD")
		return
	}
	req.ActionByName = user.Name
	req.ActionById = user.Id
	resp, err := financeClient.PaymentAdd(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// PaymentReturn godoc
// @Summary ADMIN , CEO
// @Description Return a payment for a student
// @Tags payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body pb.PaymentReturnRequest true "Payment Return Request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Router /api/finance/payment/student/return [post]
func PaymentReturn(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	user, err := utils.GetUserFromContext(ctx)
	req := pb.PaymentReturnRequest{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.ActionByName = user.Name
	req.ActionById = user.Id
	resp, err := financeClient.PaymentReturn(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// PaymentUpdate godoc
// @Summary ADMIN , CEO
// @Description Update a payment for a student
// @Tags payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body pb.PaymentUpdateRequest true "Payment Update Request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Router /api/finance/payment/student/update [patch]
func PaymentUpdate(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	user, err := utils.GetUserFromContext(ctx)
	req := pb.PaymentUpdateRequest{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	amount, err := decimal.NewFromString(req.Debit)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if amount.LessThan(decimal.NewFromInt(10000)) {
		utils.RespondError(ctx, http.StatusBadRequest, "min 10 000 USD")
		return
	}
	req.ActionByName = user.Name
	req.ActionById = user.Id
	resp, err := financeClient.PaymentUpdate(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetMonthlyStatusPayment godoc
// @Summary ADMIN , CEO
// @Description Get the monthly payment status for a student.
// @Tags payments
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {object} pb.GetMonthlyStatusResponse
// @Failure 400 {object} utils.AbsResponse "Invalid request parameters or failed retrieval"
// @Security Bearer
// @Router /api/finance/payment/student/get-monthly-status/{studentId} [get]
func GetMonthlyStatusPayment(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	studentId := ctx.Param("studentId")
	resp, err := financeClient.GetMonthlyStatusPayment(ctxR, studentId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetAllPayments godoc
// @Summary ADMIN , CEO
// @Description Get all payment records for a student within a specified month.
// @Tags payments
// @Produce json
// @Param studentId path string true "Student ID"
// @Param month path string true "Month (format: YYYY-MM)"
// @Success 200 {object} pb.GetAllPaymentsByMonthResponse
// @Failure 400 {object} utils.AbsResponse "Invalid request parameters or failed retrieval"
// @Security  Bearer
// @Router /api/finance/payment/get-all-payments/{studentId}/{month} [get]
func GetAllPayments(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	studentId := ctx.Param("studentId")
	month := ctx.Param("month")
	resp, err := financeClient.GetAllPayments(ctxR, month, studentId)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetSalaryAllTeacher godoc
// @Summary CEO
// @Description Retrieves the salary information for all teachers
// @Tags salary
// @Produce json
// @Success 200 {object} pb.GetTeachersSalaryRequest
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Security Bearer
// @Router /api/finance/salary/teacher-all [get]
func GetSalaryAllTeacher(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := financeClient.GetSalaryAllTeacher(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// AddSalaryTeacher godoc
// @Summary CEO
// @Description Adds a new salary entry for a specific teacher
// @Tags salary
// @Accept json
// @Produce json
// @Param request body pb.CreateTeacherSalaryRequest true "Salary details"
// @Success 200 {object} utils.AbsResponse "Salary added successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid request body"
// @Failure 409 {object} utils.AbsResponse "Conflict error"
// @Security Bearer
// @Router /api/finance/salary/teacher-add [post]
func AddSalaryTeacher(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req := pb.CreateTeacherSalaryRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := financeClient.AddSalaryTeacher(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteTeacherSalary godoc
// @Summary CEO
// @Description Deletes a salary entry for a specific teacher by ID
// @Tags salary
// @Produce json
// @Param teacherID path string true "Teacher ID"
// @Success 200 {object} utils.AbsResponse "Salary deleted successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict error"
// @Security Bearer
// @Router /api/finance/salary/delete/{teacherID} [delete]
func DeleteTeacherSalary(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	teacherId := ctx.Param("teacherID")
	resp, err := financeClient.DeleteTeacherSalary(ctxR, teacherId)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllTakeOffPayment handles the HTTP request to retrieve all take-off payments.
// @Summary ADMIN , CEO
// @Description Retrieve all take-off payments within a specified date range.
// @Tags payments
// @Accept json
// @Produce json
// @Param from path string true "Start date (format: YYYY-MM-DD)"
// @Param to path string true "End date (format: YYYY-MM-DD)"
// @Success 200 {object} pb.GetAllPaymentTakeOffResponse "List of take-off payments"
// @Failure 409 {object} utils.AbsResponse "Error message"
// @Failure 500 {object} utils.AbsResponse "Server error"
// @Security Bearer
// @Router /api/finance/payment/payment-take-off/{from}/{to} [get]
func GetAllTakeOffPayment(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	from := ctx.Param("from")
	to := ctx.Param("to")
	resp, err := financeClient.GetAllTakeOfPayment(from, to, ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetPaymentTakeOffChart handles the HTTP request to retrieve a chart of take-off payments.
// @Summary ADMIN , CEO
// @Description Retrieve a chart representation of take-off payments within a specified date range.
// @Tags payments
// @Accept json
// @Produce json
// @Param from path string true "Start date (format: YYYY-MM-DD)"
// @Param to path string true "End date (format: YYYY-MM-DD)"
// @Success 200 {object} pb.GetAllPaymentTakeOffChartResponse "Chart data of take-off payments"
// @Failure 409 {object} utils.AbsResponse "Error message"
// @Failure 500 {object} utils.AbsResponse "Server error"
// @Security Bearer
// @Router /api/finance/payment-takeoff-chart/{from}/{to} [get]
func GetPaymentTakeOffChart(ctx *gin.Context) {
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	from := ctx.Param("from")
	to := ctx.Param("to")
	resp, err := financeClient.GetPaymentTakeOffChart(from, to, ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
