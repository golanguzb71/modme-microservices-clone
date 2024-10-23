package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{client: client}, nil
}

func (c *UserClient) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.AbsResponse, error) {
	return c.client.CreateUser(ctx, req)
}

func (c *UserClient) GetTeachers(ctx context.Context, isDeleted bool) (*pb.GetTeachersResponse, error) {
	return c.client.GetTeachers(ctx, &pb.GetTeachersRequest{IsDeleted: isDeleted})
}
