package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

type FinanceClient struct {
	discountClient      pb.DiscountServiceClient
	categoryClient      pb.CategoryServiceClient
	expenseClient       pb.ExpenseServiceClient
	paymentClient       pb.PaymentServiceClient
	teacherSalaryClient pb.TeacherSalaryServiceClient
}

func (fc *FinanceClient) GetDiscountsInformationByGroupId(ctx context.Context, groupId string) (*pb.GetInformationDiscountResponse, error) {
	return fc.discountClient.GetAllInformationDiscount(ctx, &pb.GetInformationDiscountRequest{GroupId: groupId})
}
func (fc *FinanceClient) CreateDiscount(ctx context.Context, req *pb.AbsDiscountRequest) (*pb.AbsResponse, error) {
	return fc.discountClient.CreateDiscount(ctx, req)
}
func (fc *FinanceClient) DeleteDiscount(ctx context.Context, groupId string, studentId string) (*pb.AbsResponse, error) {
	return fc.discountClient.DeleteDiscount(ctx, &pb.AbsDiscountRequest{
		GroupId:       groupId,
		StudentId:     studentId,
		DiscountPrice: "",
		Comment:       "",
	})
}
func (fc *FinanceClient) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.AbsResponse, error) {
	return fc.categoryClient.CreateCategory(ctx, req)
}
func (fc *FinanceClient) DeleteCategory(ctx context.Context, req string) (*pb.AbsResponse, error) {
	return fc.categoryClient.DeleteCategory(ctx, &pb.DeleteAbsRequest{Id: req})
}
func (fc *FinanceClient) GetAllCategories(ctx context.Context) (*pb.GetAllCategoryRequest, error) {
	return fc.categoryClient.GetAllCategory(ctx, &emptypb.Empty{})
}

func (fc *FinanceClient) CreateExpense(ctx context.Context, req *pb.CreateExpenseRequest) (*pb.AbsResponse, error) {
	return fc.expenseClient.CreateExpense(ctx, req)
}

func (fc *FinanceClient) DeleteExpense(ctx context.Context, id string) (*pb.AbsResponse, error) {
	return fc.expenseClient.DeleteExpense(ctx, &pb.DeleteAbsRequest{Id: id})
}

func (fc *FinanceClient) GetAllInformation(ctx context.Context, id string, idType string, page int64, size int64, from string, to string) (*pb.GetAllExpenseResponse, error) {
	return fc.expenseClient.GetAllExpense(ctx, &pb.GetAllExpenseRequest{
		From: from,
		To:   to,
		Type: idType,
		Id:   id,
		PageReq: &pb.PageRequest{
			Page: int32(page),
			Size: int32(size),
		},
	})
}

func (fc *FinanceClient) GetHistoryDiscount(id string, ctx context.Context) (*pb.GetHistoryDiscountResponse, error) {
	return fc.discountClient.GetHistoryDiscount(ctx, &pb.GetHistoryDiscountRequest{StudentId: id})
}

func (fc *FinanceClient) GetExpenseChartDiagram(from string, to string, ctx context.Context) (*pb.GetAllExpenseDiagramResponse, error) {
	return fc.expenseClient.GetAllExpenseDiagram(ctx, &pb.GetAllExpenseDiagramRequest{
		From: from,
		To:   to,
	})
}

func (fc *FinanceClient) PaymentAdd(req *pb.PaymentAddRequest) (*pb.AbsResponse, error) {
	return fc.paymentClient.PaymentAdd(context.TODO(), req)
}

func (fc *FinanceClient) PaymentReturn(ctx context.Context, req *pb.PaymentReturnRequest) (*pb.AbsResponse, error) {
	return fc.paymentClient.PaymentReturn(ctx, req)
}

func (fc *FinanceClient) PaymentUpdate(ctx context.Context, p *pb.PaymentUpdateRequest) (*pb.AbsResponse, error) {
	return fc.paymentClient.PaymentUpdate(ctx, p)
}

func (fc *FinanceClient) GetMonthlyStatusPayment(ctx context.Context, studentId string) (*pb.GetMonthlyStatusResponse, error) {
	return fc.paymentClient.GetMonthlyStatus(ctx, &pb.GetMonthlyStatusRequest{UserId: studentId})
}

func (fc *FinanceClient) GetAllPayments(ctx context.Context, month string, studentId string) (*pb.GetAllPaymentsByMonthResponse, error) {
	resp, err := fc.paymentClient.GetAllPaymentsByMonth(ctx, &pb.GetAllPaymentsByMonthRequest{
		UserId: studentId,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}

	for _, payment := range resp.Payments {
		amountFloat, parseErr := strconv.ParseFloat(payment.Amount, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse amount '%s': %v", payment.Amount, parseErr)
		}
		payment.Amount = fmt.Sprintf("%.0f", amountFloat)
	}

	return resp, nil
}

func (fc *FinanceClient) GetSalaryAllTeacher(ctx context.Context) (*pb.GetTeachersSalaryRequest, error) {
	return fc.teacherSalaryClient.GetTeacherSalary(ctx, &emptypb.Empty{})
}

func (fc *FinanceClient) AddSalaryTeacher(ctx context.Context, req *pb.CreateTeacherSalaryRequest) (*pb.AbsResponse, error) {
	return fc.teacherSalaryClient.CreateTeacherSalary(ctx, req)
}

func (fc *FinanceClient) DeleteTeacherSalary(ctx context.Context, teacherId string) (*pb.AbsResponse, error) {
	return fc.teacherSalaryClient.DeleteTeacherSalary(ctx, &pb.DeleteTeacherSalaryRequest{TeacherId: teacherId})
}

func (fc *FinanceClient) GetAllTakeOfPayment(from string, to string, ctx context.Context) (*pb.GetAllPaymentTakeOffResponse, error) {
	return fc.paymentClient.GetAllPaymentTakeOff(ctx, &pb.GetAllPaymentTakeOffRequest{
		From: from,
		To:   to,
	})
}

func (fc *FinanceClient) GetPaymentTakeOffChart(from string, to string, ctx context.Context) (*pb.GetAllPaymentTakeOffChartResponse, error) {
	return fc.paymentClient.GetAllPaymentTakeOffChart(ctx, &pb.GetAllPaymentTakeOffRequest{
		From: from,
		To:   to,
	})
}

func (fc *FinanceClient) GetAllStudentPayment(from string, to string, ctx context.Context) (*pb.GetAllStudentPaymentsResponse, error) {
	return fc.paymentClient.GetAllStudentPayments(ctx, &pb.GetAllStudentPaymentsRequest{
		From: from,
		To:   to,
	})
}

func (fc *FinanceClient) GetAllPaymentsStudent(from string, to string, ctx context.Context) (*pb.GetAllStudentPaymentsChartResponse, error) {
	return fc.paymentClient.GetAllStudentPaymentsChart(ctx, &pb.GetAllStudentPaymentsRequest{
		From: from,
		To:   to,
	})
}

func (fc *FinanceClient) GetAllDebtsInformation(ctx context.Context, page, size, from, to string) (*pb.GetAllDebtsInformationResponse, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}
	return fc.paymentClient.GetAllDebtsInformation(ctx, &pb.GetAllDebtsRequest{
		PageParam: &pb.PageRequest{
			Page: int32(pageInt),
			Size: int32(sizeInt),
		},
		From: from,
		To:   to,
	})
}

func (fc *FinanceClient) GetCommonFinanceInformation(ctx context.Context) (int, int) {
	response, err := fc.paymentClient.GetCommonFinanceInformation(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, 0
	}
	return int(response.DebtorsCount), int(response.PayInCurrentMonth)
}

func (fc *FinanceClient) GetChartIncome(ctx context.Context, from string, to string) (*pb.GetCommonInformationResponse, error) {
	_, err := fc.paymentClient.GetIncomeChart(ctx, &pb.GetIncomeChartRequest{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (fc *FinanceClient) GetTableGroups(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func NewFinanceClient(addr string) (*FinanceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	discountClient := pb.NewDiscountServiceClient(conn)
	categoryClient := pb.NewCategoryServiceClient(conn)
	expenseClient := pb.NewExpenseServiceClient(conn)
	paymentClient := pb.NewPaymentServiceClient(conn)
	teacherClient := pb.NewTeacherSalaryServiceClient(conn)
	return &FinanceClient{discountClient: discountClient, categoryClient: categoryClient, expenseClient: expenseClient, paymentClient: paymentClient, teacherSalaryClient: teacherClient}, nil
}
