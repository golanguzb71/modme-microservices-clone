package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FinanceClient struct {
	discountClient pb.DiscountServiceClient
	categoryClient pb.CategoryServiceClient
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

func NewFinanceClient(addr string) (*FinanceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	discountClient := pb.NewDiscountServiceClient(conn)
	categoryClient := pb.NewCategoryServiceClient(conn)
	return &FinanceClient{discountClient: discountClient, categoryClient: categoryClient}, nil
}
