package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/internal/utils"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type ExpenseService struct {
	pb.UnimplementedExpenseServiceServer
	repo *repository.ExpenseRepository
}

func (e *ExpenseService) CreateExpense(ctx context.Context, req *pb.CreateExpenseRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := e.repo.CreateExpense(req.Title, req.GivenDate, req.ExpenseType, req.CategoryId, req.UserId, req.Sum, req.CreatedById, req.PaymentMethod); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusCreated,
		Message: "expense created",
	}, nil
}
func (e *ExpenseService) DeleteExpense(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := e.repo.DeleteExpense(req.Id); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "expense deleted",
	}, nil
}
func (e *ExpenseService) GetAllExpense(ctx context.Context, req *pb.GetAllExpenseRequest) (*pb.GetAllExpenseResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return e.repo.GetAllExpense(req.PageReq.Page, req.PageReq.Size, req.From, req.To, req.Type, req.Id)
}
func (e *ExpenseService) GetAllExpenseDiagram(ctx context.Context, req *pb.GetAllExpenseDiagramRequest) (*pb.GetAllExpenseDiagramResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return e.repo.GetExpenseDiagram(req.To, req.From)
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}
