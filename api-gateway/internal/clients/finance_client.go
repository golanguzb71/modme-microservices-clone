package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
)

type FinanceClient struct {
	discountClient pb.DiscountServiceClient
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

func NewFinanceClient(addr string) (*FinanceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	discountClient := pb.NewDiscountServiceClient(conn)
	return &FinanceClient{discountClient: discountClient}, nil
}
