package service

import (
	"context"
	"finance-service/internal/repository"
	"finance-service/internal/utils"
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
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if req.Type == "ADD" {
		if err := ps.repo.AddPayment(ctx, companyId, req.Date, req.Sum, req.Method, req.Comment, req.UserId, req.ActionByName, req.ActionById, req.GroupId, false); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "payment added",
		}, nil
	} else if req.Type == "TAKE_OFF" {
		if err := ps.repo.TakeOffPayment(ctx, companyId, req.Date, req.Sum, req.Method, req.Comment, req.UserId, req.ActionByName, req.ActionById, req.GroupId); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}
		return &pb.AbsResponse{
			Status:  http.StatusCreated,
			Message: "payment take_off successfully",
		}, nil
	} else if req.Type == "REFUND" {
		if err := ps.repo.AddPayment(ctx, companyId, req.Date, req.Sum, req.Method, req.Comment, req.UserId, req.ActionByName, req.ActionById, req.GroupId, true); err != nil {
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
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.PaymentReturn(ctx, companyId, req.PaymentId, req.ActionByName, req.ActionById)
}
func (ps *PaymentService) PaymentUpdate(ctx context.Context, req *pb.PaymentUpdateRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.PaymentUpdate(ctx, companyId, req.PaymentId, req.Date, req.Method, req.UserId, req.Comment, req.Debit, req.ActionByName, req.ActionById, req.GroupId)
}
func (ps *PaymentService) GetMonthlyStatus(ctx context.Context, req *pb.GetMonthlyStatusRequest) (*pb.GetMonthlyStatusResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetMonthlyStatus(ctx, companyId, req.UserId)
}
func (ps *PaymentService) GetAllPaymentsByMonth(ctx context.Context, req *pb.GetAllPaymentsByMonthRequest) (*pb.GetAllPaymentsByMonthResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetAllPaymentsByMonth(ctx, companyId, req.Month, req.UserId)
}

func (ps *PaymentService) GetAllPaymentTakeOff(ctx context.Context, req *pb.GetAllPaymentTakeOffRequest) (*pb.GetAllPaymentTakeOffResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetAllPaymentTakeOff(ctx, companyId, req.From, req.To)
}
func (ps *PaymentService) GetAllPaymentTakeOffChart(ctx context.Context, req *pb.GetAllPaymentTakeOffRequest) (*pb.GetAllPaymentTakeOffChartResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetAllPaymentTakeOffChart(ctx, companyId, req.From, req.To)
}

func (ps *PaymentService) GetAllStudentPayments(ctx context.Context, req *pb.GetAllStudentPaymentsRequest) (*pb.GetAllStudentPaymentsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetAllStudentPayments(ctx, companyId, req.From, req.To)
}

func (ps *PaymentService) GetAllStudentPaymentsChart(ctx context.Context, req *pb.GetAllStudentPaymentsRequest) (*pb.GetAllStudentPaymentsChartResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetAllStudentPaymentsChart(ctx, companyId, req.From, req.To)
}

func (ps *PaymentService) GetAllDebtsInformation(ctx context.Context, req *pb.GetAllDebtsRequest) (*pb.GetAllDebtsInformationResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetAllDebtsInformation(ctx, companyId, req.From, req.To, req.PageParam.Page, req.PageParam.Size)
}
func (ps *PaymentService) GetCommonFinanceInformation(ctx context.Context, req *emptypb.Empty) (*pb.GetCommonInformationResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetCommonFinanceInformation(ctx, companyId)
}

func (ps *PaymentService) GetIncomeChart(ctx context.Context, req *pb.GetIncomeChartRequest) (*pb.GetIncomeChartResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return ps.repo.GetIncomeChart(ctx, companyId, req.From, req.To)
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}
