package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EducationClient struct {
	roomClient       pb.RoomServiceClient
	courseClient     pb.CourseServiceClient
	groupClient      pb.GroupServiceClient
	attendanceClient pb.AttendanceServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	roomClient := pb.NewRoomServiceClient(conn)
	courseClient := pb.NewCourseServiceClient(conn)
	groupClient := pb.NewGroupServiceClient(conn)
	attendanceClient := pb.NewAttendanceServiceClient(conn)
	return &EducationClient{roomClient: roomClient, courseClient: courseClient, groupClient: groupClient, attendanceClient: attendanceClient}, nil
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

func (lc *EducationClient) GetCourseById(ctx context.Context, id string) (*pb.GetCourseByIdResponse, error) {
	return lc.courseClient.GetCourseById(ctx, &pb.GetCourseByIdRequest{Id: id})
}

func (lc *EducationClient) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.AbsResponse, error) {
	return lc.groupClient.CreateGroup(ctx, req)
}

func (lc *EducationClient) UpdateGroup(ctx context.Context, req *pb.GetUpdateGroupAbs) (*pb.AbsResponse, error) {
	return lc.groupClient.UpdateGroup(ctx, req)
}

func (lc *EducationClient) DeleteGroup(ctx context.Context, id string) (*pb.AbsResponse, error) {
	return lc.groupClient.DeleteGroup(ctx, &pb.DeleteAbsRequest{Id: id})
}

func (lc *EducationClient) GetAllGroup(ctx context.Context, isArchived bool, page, size int32) (*pb.GetGroupsResponse, error) {
	return lc.groupClient.GetGroups(ctx, &pb.GetGroupsRequest{
		IsArchived: isArchived,
		Page: &pb.PageRequest{
			Page: page,
			Size: size,
		},
	})
}

func (lc *EducationClient) GetGroupById(ctx context.Context, id string) (*pb.GetGroupAbsResponse, error) {
	return lc.groupClient.GetGroupById(ctx, &pb.GetGroupByIdRequest{Id: id})
}

func (lc *EducationClient) GetAttendanceByGroup(ctx context.Context, req *pb.GetAttendanceRequest) (*pb.GetAttendanceResponse, error) {
	return lc.attendanceClient.GetAttendance(ctx, req)
}

func (lc *EducationClient) SetAttendanceByGroup(ctx context.Context, req *pb.SetAttendanceRequest) (*pb.AbsResponse, error) {
	return lc.attendanceClient.SetAttendance(ctx, req)
}
