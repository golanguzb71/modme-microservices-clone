package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EducationClient struct {
	roomClient pb.RoomServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	roomClient := pb.NewRoomServiceClient(conn)
	return &EducationClient{roomClient: roomClient}, nil
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
