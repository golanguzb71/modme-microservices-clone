package clients

import (
	"context"
	"education-service/proto/pb"
	"google.golang.org/grpc"
)

type FinanceClient struct {
	discountClient pb.DiscountServiceClient
}

func NewFinanceClient(addr string) (*FinanceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	discountClient := pb.NewDiscountServiceClient(conn)
	return &FinanceClient{discountClient: discountClient}, nil
}

func (fc *FinanceClient) GetDiscountByStudentId(ctx context.Context, studentId string) *string {
	resp, err := fc.discountClient.GetDiscountByStudentId(ctx, &pb.GetDiscountByStudentIdRequest{StudentId: studentId})
	if err != nil {
		return nil
	}
	if !resp.IsHave {
		return nil
	}
	return &resp.Amount
}
