package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoomService struct {
	pb.UnimplementedRoomServiceServer
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.AbsResponse, error) {
	err := s.repo.CreateRoom(req.Name, req.Capacity)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Room created successfully"}, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, req *pb.AbsRoom) (*pb.AbsResponse, error) {
	err := s.repo.UpdateRoom(&req.Id, &req.Name, &req.Capacity)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Room updated successfully"}, nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	err := s.repo.DeleteRoom(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{Status: 200, Message: "Expectation deleted successfully"}, nil
}

func (s *RoomService) GetRooms(ctx context.Context, req *emptypb.Empty) (*pb.GetUpdateRoomAbs, error) {
	return s.repo.GetRoom()
}
