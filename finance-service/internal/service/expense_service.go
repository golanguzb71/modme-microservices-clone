package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/proto/pb"
	"net/http"
)

type ExpenseService struct {
	pb.UnimplementedExpenseServiceServer
	repo *repository.ExpenseRepository
}

func (e *ExpenseService) CreateExpense(ctx context.Context, req *pb.CreateExpenseRequest) (*pb.AbsResponse, error) {
	if err := e.repo.CreateExpense(req.Title, req.GivenDate, req.ExpenseType, req.CategoryId, req.UserId, req.Sum, req.CreatedById, req.PaymentMethod); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusCreated,
		Message: "expense created",
	}, nil
}
func (e *ExpenseService) DeleteExpense(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	if err := e.repo.DeleteExpense(req.Id); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "expense deleted",
	}, nil
}
func (e *ExpenseService) GetAllExpense(ctx context.Context, req *pb.GetAllExpenseRequest) (*pb.GetAllExpenseResponse, error) {
	return e.repo.GetAllExpense(req.PageReq.Page, req.PageReq.Size, req.From, req.To, req.ByCategory)
}
func (e *ExpenseService) GetAllExpenseDiagram(ctx context.Context, req *pb.GetAllExpenseDiagramRequest) (*pb.GetAllExpenseDiagramResponse, error) {
	return e.repo.GetExpenseDiagram(req.To, req.From)
}
