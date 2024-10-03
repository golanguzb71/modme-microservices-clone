package client

import (
	"api-gateway/grpc/proto/pb"
	"google.golang.org/grpc"
)

type BusinessClient struct {
	client pb.BusinessServiceClient
}

func NewBusinessClient(addr string) (*BusinessClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewBusinessServiceClient(conn)
	return &BusinessClient{client: client}, nil
}
