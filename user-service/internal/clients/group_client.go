package clients

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"user-service/proto/pb"
)

type GroupClient struct {
	client pb.GroupServiceClient
}

func NewGroupClient(addr string) (*GroupClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GroupService: %w", err)
	}
	client := pb.NewGroupServiceClient(conn)
	return &GroupClient{client: client}, nil
}

func (gc *GroupClient) GetGroupsByTeacherId(ctx context.Context, teacherId string, isArchived bool) (int, error) {
	if gc.client == nil {
		return 0, fmt.Errorf("GroupService client is not initialized")
	}
	resp, err := gc.client.GetGroupsByTeacherId(ctx, &pb.GetGroupsByTeacherIdRequest{
		TeacherId:  teacherId,
		IsArchived: isArchived,
	})
	if err != nil {
		return 0, err
	}
	return len(resp.Groups), nil
}
