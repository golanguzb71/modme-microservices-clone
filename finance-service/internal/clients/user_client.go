package clients

import (
	"context"
	"finance-service/proto/pb"
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

func (c *UserClient) GetUserById(ctx context.Context, teacherId string) (*pb.GetUserByIdResponse, error) {
	return c.client.GetUserById(ctx, &pb.UserAbsRequest{UserId: teacherId})
}
