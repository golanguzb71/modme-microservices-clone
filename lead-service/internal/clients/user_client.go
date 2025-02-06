package clients

import (
	"context"
	"google.golang.org/grpc"
	"lid-service/proto/pb"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(addr string) *UserClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil
	}
	client := pb.NewUserServiceClient(conn)
	return &UserClient{client: client}
}

func (uc *UserClient) GetUserById(ctx context.Context, id string) (*pb.GetUserByIdResponse, error) {
	return uc.client.GetUserById(ctx, &pb.UserAbsRequest{UserId: id})
}
