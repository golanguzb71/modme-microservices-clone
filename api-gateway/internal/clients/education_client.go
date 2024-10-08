package client

import (
	"api-gateway/grpc/proto/pb"
	"google.golang.org/grpc"
)

type EducationClient struct {
	client pb.EducationServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewEducationServiceClient(conn)
	return &EducationClient{client: client}, nil
}
