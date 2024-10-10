package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EducationClient struct {
	roomClient   pb.RoomServiceClient
	courseClient pb.CourseServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	roomClient := pb.NewRoomServiceClient(conn)
	courseClient := pb.NewCourseServiceClient(conn)
	return &EducationClient{roomClient: roomClient, courseClient: courseClient}, nil
}

// Education Service method client

func (lc *EducationClient) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.AbsResponse, error) {
	return lc.roomClient.CreateRoom(ctx, req)
}

func (lc *EducationClient) UpdateRoom(ctx context.Context, req *pb.AbsRoom) (*pb.AbsResponse, error) {
	return lc.roomClient.UpdateRoom(ctx, req)
}

func (lc *EducationClient) DeleteRoom(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := pb.DeleteAbsRequest{
		Id: id,
	}
	return lc.roomClient.DeleteRoom(ctx, &req)
}

func (lc *EducationClient) GetRoom(ctx context.Context) (*pb.GetUpdateRoomAbs, error) {
	return lc.roomClient.GetRooms(ctx, &emptypb.Empty{})
}

func (lc *EducationClient) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.AbsResponse, error) {
	return lc.courseClient.CreateCourse(ctx, req)
}

func (lc *EducationClient) UpdateCourse(ctx context.Context, req *pb.AbsCourse) (*pb.AbsResponse, error) {
	return lc.courseClient.UpdateCourse(ctx, req)
}

func (lc *EducationClient) DeleteCourse(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := pb.DeleteAbsRequest{
		Id: id,
	}
	return lc.courseClient.DeleteCourse(ctx, &req)
}

func (lc *EducationClient) GetCourse(ctx context.Context) (*pb.GetUpdateCourseAbs, error) {
	return lc.courseClient.GetCourses(ctx, &emptypb.Empty{})
}
