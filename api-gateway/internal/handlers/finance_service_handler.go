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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

func GetChartDiagram(ctx *gin.Context) {

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
	ctxR, cancel := context.WithTimeout(context.Background(), time.Second*5)
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