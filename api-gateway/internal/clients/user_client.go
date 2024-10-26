package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
	"time"
)

type UserClient struct {
	client     pb.UserServiceClient
	authClient pb.AuthServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	authClient := pb.NewAuthServiceClient(conn)
	return &UserClient{client: client, authClient: authClient}, nil
}

func (c *UserClient) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.AbsResponse, error) {
	return c.client.CreateUser(ctx, req)
}

func (c *UserClient) GetTeachers(ctx context.Context, isDeleted bool) (*pb.GetTeachersResponse, error) {
	return c.client.GetTeachers(ctx, &pb.GetTeachersRequest{IsDeleted: isDeleted})
}

func (c *UserClient) GetUserById(ctx context.Context, teacherId string) (*pb.GetUserByIdResponse, error) {
	return c.client.GetUserById(ctx, &pb.UserAbsRequest{UserId: teacherId})
}

func (c *UserClient) UpdateUserById(ctx context.Context, req *pb.UpdateUserRequest) (*pb.AbsResponse, error) {
	return c.client.UpdateUserById(ctx, req)
}

func (c *UserClient) DeleteUserById(ctx context.Context, userId string) (*pb.AbsResponse, error) {
	return c.client.DeleteUserById(ctx, &pb.UserAbsRequest{UserId: userId})
}

func (c *UserClient) GetAllEmployee(ctx context.Context, isArchived bool) (*pb.GetAllEmployeeResponse, error) {
	return c.client.GetAllEmployee(ctx, &pb.GetAllEmployeeRequest{IsArchived: isArchived})
}

func (c *UserClient) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return c.authClient.Login(ctx, request)
}

func (c *UserClient) ValidateToken(token string, requiredRoles []string) (*pb.GetUserByIdResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	req := &pb.ValidateTokenRequest{
		Token:         token,
		RequiredRoles: requiredRoles,
	}
	resp, err := c.authClient.ValidateToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
