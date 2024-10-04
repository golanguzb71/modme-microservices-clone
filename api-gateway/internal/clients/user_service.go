package client

import (
	"api-gateway/grpc/proto/pb"
	"google.golang.org/grpc"
)

type UserClient struct {
	client pb.UserServiceClient
}

//func (c UserClient) ValidateToken(token string, roles []string) (*pb.User, error) {
//
//}
//

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{client: client}, nil
}
