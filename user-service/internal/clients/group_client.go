package clients

import (
	"context"
	"google.golang.org/grpc"
	"user-service/proto/pb"
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

func (gc *GroupClient) GetGroupsByTeacherId(teacherId string, isArchived bool) (int, error) {
	resp, err := gc.client.GetGroupsByTeacherId(context.TODO(), &pb.GetGroupsByTeacherIdRequest{TeacherId: teacherId, IsArchived: isArchived})
	if err != nil {
		return 0, err
	}
	return len(resp.Groups), nil
}
