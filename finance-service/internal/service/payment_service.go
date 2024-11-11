package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	repo *repository.PaymentRepository
}

func (ps *PaymentService) PaymentAdd(ctx context.Context, req *pb.PaymentAddRequest) (*pb.AbsResponse, error) {
	if req.Type == "ADD" {
		if err := ps.repo.AddPayment(req.Date, req.Sum, req.Method, req.Comment, req.UserId, req.ActionByName, req.ActionById, req.GroupId, false); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "payment added",
		}, nil
	} else if req.Type == "TAKE_OFF" {
		if err := ps.repo.TakeOffPayment(req.Date, req.Sum, req.Method, req.Comment, req.UserId, req.ActionByName, req.ActionById, req.GroupId); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "payment take_off successfully",
		}, nil
	} else if req.Type == "REFUND" {
		if err := ps.repo.AddPayment(req.Date, req.Sum, req.Method, req.Comment, req.UserId, req.ActionByName, req.ActionById, req.GroupId, true); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "payment refund",
		}, nil
	}
	return nil, status.Errorf(codes.Aborted, "invalid request type")
}
func (ps *PaymentService) PaymentReturn(ctx context.Context, req *pb.PaymentReturnRequest) (*pb.AbsResponse, error) {
	return ps.repo.PaymentReturn(req.PaymentId, req.ActionByName, req.ActionById)
}
func (ps *PaymentService) PaymentUpdate(ctx context.Context, req *pb.PaymentUpdateRequest) (*pb.AbsResponse, error) {
	return ps.repo.PaymentUpdate(req.PaymentId, req.Date, req.Method, req.UserId, req.Comment, req.Debit, req.ActionByName, req.ActionById, req.GroupId)
}
func (ps *PaymentService) GetMonthlyStatus(ctx context.Context, req *pb.GetMonthlyStatusRequest) (*pb.GetMonthlyStatusResponse, error) {
	return ps.repo.GetMonthlyStatus(req.UserId)
}
func (ps *PaymentService) GetAllPaymentsByMonth(ctx context.Context, req *pb.GetAllPaymentsByMonthRequest) (*pb.GetAllPaymentsByMonthResponse, error) {
	return ps.repo.GetAllPaymentsByMonth(req.Month, req.UserId)
}

func (ps *PaymentService) GetAllPaymentTakeOff(ctx context.Context, req *pb.GetAllPaymentTakeOffRequest) (*pb.GetAllPaymentTakeOffResponse, error) {
	return ps.repo.GetAllPaymentTakeOff(req.From, req.To)
}
func (ps *PaymentService) GetAllPaymentTakeOffChart(ctx context.Context, req *pb.GetAllPaymentTakeOffRequest) (*pb.GetAllPaymentTakeOffChartResponse, error) {
	return ps.repo.GetAllPaymentTakeOffChart(req.From, req.To)
}

func (ps *PaymentService) GetAllStudentPayments(ctx context.Context, req *pb.GetAllStudentPaymentsRequest) (*pb.GetAllStudentPaymentsResponse, error) {
	return ps.repo.GetAllStudentPayments(req.From, req.To)
}

func (ps *PaymentService) GetAllStudentPaymentsChart(ctx context.Context, req *pb.GetAllStudentPaymentsRequest) (*pb.GetAllStudentPaymentsChartResponse, error) {
	return ps.repo.GetAllStudentPaymentsChart(req.From, req.To)
}

func (ps *PaymentService) GetAllDebtsInformation(ctx context.Context, req *pb.GetAllDebtsRequest) (*pb.GetAllDebtsInformationResponse, error) {
	return ps.repo.GetAllDebtsInformation(req.From, req.To, req.PageParam.Page, req.PageParam.Size)
}
func (ps *PaymentService) GetCommonFinanceInformation(ctx context.Context, req *emptypb.Empty) (*pb.GetCommonInformationResponse, error) {
	return ps.repo.GetCommonFinanceInformation()
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}
