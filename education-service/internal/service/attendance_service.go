package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
	"errors"
	"time"
)

type AttendanceService struct {
	pb.UnimplementedAttendanceServiceServer
	attendanceRepo *repository.AttendanceRepository
}

func NewAttendanceService(repo *repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: repo,
	}
}

func (s *AttendanceService) GetAttendance(ctx context.Context, req *pb.GetAttendanceRequest) (*pb.GetAttendanceResponse, error) {
	if req.GroupId == "" {
		return nil, errors.New("group ID is required")
	}
	fromDate, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		return nil, errors.New("invalid 'from' date format")
	}
	tillDate, err := time.Parse("2006-01-02", req.Till)
	if err != nil {
		return nil, errors.New("invalid 'till' date format")
	}
	if tillDate.Before(fromDate) {
		return nil, errors.New("'till' date must be after 'from' date")
	}
	return s.attendanceRepo.GetAttendanceByGroupAndDateRange(ctx, req.GroupId, fromDate, tillDate)
}

func (s *AttendanceService) SetAttendance(ctx context.Context, req *pb.SetAttendanceRequest) (*pb.AbsResponse, error) {
	if req.GroupId == "" || req.StudentId == "" || req.TeacherId == "" {
		return nil, errors.New("group ID, student ID, and teacher ID are required")
	}

	attendDate, err := time.Parse("2006-01-02", req.AttendDate)
	if err != nil {
		return nil, errors.New("invalid attendance date format")
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	if attendDate.After(today) {
		return nil, errors.New("attendance date cannot be in the future")
	}

	validDay, err := s.attendanceRepo.IsValidGroupDay(ctx, req.GroupId, today)
	if err != nil {
		return nil, err
	}
	if !validDay {
		return nil, errors.New("attendance cannot be created today; group is not active")
	}

	cutoffTime := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	if attendDate.Equal(today.AddDate(0, 0, -1)) && now.After(cutoffTime) {
		return nil, errors.New("attendance cannot be set for yesterday after 12 PM")
	}

	err = s.attendanceRepo.CreateAttendance(ctx, req.GroupId, req.StudentId, req.TeacherId, req.AttendDate, req.Status)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "Attendance successfully created",
	}, nil
}
