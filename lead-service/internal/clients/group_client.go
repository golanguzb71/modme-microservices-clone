package clients

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"lid-service/proto/pb"
)

type GroupClient struct {
	client       pb.GroupServiceClient
	courseClient pb.CourseServiceClient
}

func NewGroupClient(addr string) *GroupClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := pb.NewGroupServiceClient(conn)
	courseClient := pb.NewCourseServiceClient(conn)
	return &GroupClient{client: client, courseClient: courseClient}
}

func (gc *GroupClient) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (error, string) {
	resp, err := gc.client.CreateGroup(ctx, req)
	if err != nil {
		return err, ""
	}
	return nil, resp.Message
}

func (gc *GroupClient) GetCourse(ctx context.Context, id string) (*pb.GetCourseByIdResponse, error) {
	return gc.courseClient.GetCourseById(ctx, &pb.GetCourseByIdRequest{Id: id})
}
