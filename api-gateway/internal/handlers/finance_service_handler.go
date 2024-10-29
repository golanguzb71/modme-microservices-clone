package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
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

func GetAllInformation(ctx *gin.Context) {

}

func GetChartDiagram(ctx *gin.Context) {

}
