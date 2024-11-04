package clients

import (
	"context"
	"google.golang.org/grpc"
	"lid-service/proto/pb"
)

type GroupClient struct {
	client pb.GroupServiceClient
}

func NewGroupClient(addr string) *GroupClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil
	}
	client := pb.NewGroupServiceClient(conn)
	return &GroupClient{client: client}
}

func (gc *GroupClient) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (error, string) {
	resp, err := gc.client.CreateGroup(ctx, req)
	if err != nil {
		return err, ""
	}
	return nil, resp.Message
}
