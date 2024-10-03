package client

import (
	"api-gateway/grpc/proto/pb"
	"google.golang.org/grpc"
)

type AuditClient struct {
	client pb.AuditServiceClient
}

func NewAuditClient(addr string) (*AuditClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewAuditServiceClient(conn)
	return &AuditClient{client: client}, nil
}
