package clients

import (
	"context"
	"education-service/proto/pb"
	"fmt"
	"google.golang.org/grpc"
	"strconv"
)

type FinanceClient struct {
	discountClient      pb.DiscountServiceClient
	paymentClient       pb.PaymentServiceClient
	teacherSalaryClient pb.TeacherSalaryServiceClient
}

func NewFinanceClient(addr string) (*FinanceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	discountClient := pb.NewDiscountServiceClient(conn)
	paymentClient := pb.NewPaymentServiceClient(conn)
	teacherSalaryClient := pb.NewTeacherSalaryServiceClient(conn)
	return &FinanceClient{discountClient: discountClient, paymentClient: paymentClient, teacherSalaryClient: teacherSalaryClient}, nil

}

func (fc *FinanceClient) GetDiscountByStudentId(ctx context.Context, studentId, groupId string) (*float64, string) {
	resp, err := fc.discountClient.GetDiscountByStudentId(ctx, &pb.GetDiscountByStudentIdRequest{StudentId: studentId, GroupId: groupId})
	if err != nil {
		fmt.Printf("error while getting student discount %v\n", err)
		return nil, "CENTER"
	}
	if !resp.IsHave {
		fmt.Println("response is have is false", resp.IsHave)
		return nil, "CENTER"
	}
	discountAmount, err := strconv.ParseFloat(resp.Amount, 64)
	return &discountAmount, resp.DiscountOwner
}

func (fc *FinanceClient) PaymentAdd(ctx context.Context, comment, date, method, sum, userId, paymentType, actionById, actionByName, groupId string) (*pb.AbsResponse, error) {
	return fc.paymentClient.PaymentAdd(ctx, &pb.PaymentAddRequest{
		Comment:      comment,
		Date:         date,
		Method:       method,
		Sum:          sum,
		UserId:       userId,
		Type:         paymentType,
		ActionById:   actionById,
		ActionByName: actionByName,
		GroupId:      groupId,
	})
}

func (fc *FinanceClient) GetTeacherSalaryByTeacherID(ctx context.Context, teacherId string) (*pb.AbsGetTeachersSalary, error) {
	return fc.teacherSalaryClient.GetTeacherSalaryByTeacherID(ctx, &pb.DeleteTeacherSalaryRequest{TeacherId: teacherId})
}
