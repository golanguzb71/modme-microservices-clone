package client

import (
	"api-gateway/grpc/proto/pb"
	"google.golang.org/grpc"
)

type LidClient struct {
	client pb.LidServiceClient
}

func NewLidClient(addr string) (*LidClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewLidServiceClient(conn)
	return &LidClient{client: client}, nil
}
