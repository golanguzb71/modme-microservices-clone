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
// @Tags Discounts
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
// @Tags Discounts
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
// @Tags Discounts
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
