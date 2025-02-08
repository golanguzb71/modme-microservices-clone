package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
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
	return s.attendanceRepo.GetAttendanceByGroupAndDateRange(companyId, ctx, req.GroupId, fromDate, tillDate, req.WithOutdated, req.ActionRole, req.ActionId)
}
func (s *AttendanceService) SetAttendance(ctx context.Context, req *pb.SetAttendanceRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.PermissionDenied, "error while getting company from context")
	}
	if req.GroupId == "" || req.StudentId == "" || req.TeacherId == "" {
		return nil, errors.New("group ID, student ID, and teacher ID are required")
	}

	if req.ActionByRole == "CEO" || req.ActionByRole == "ADMIN" {
		if req.Status == -1 {
			err := s.attendanceRepo.DeleteAttendance(req.GroupId, req.StudentId, req.TeacherId, req.AttendDate)
			if err != nil {
				return nil, err
			}
			return &pb.AbsResponse{
				Status:  200,
				Message: "Attendance successfully deleted",
			}, nil
		} else {
			err := s.attendanceRepo.CreateAttendance(ctx, companyId, req.GroupId, req.StudentId, req.TeacherId, req.AttendDate, req.Status, req.ActionById, req.ActionByRole)
			if err != nil {
				return nil, err
			}
			return &pb.AbsResponse{
				Status:  200,
				Message: "Attendance successfully created",
			}, nil
		}
	}
	attendDate, err := time.Parse("2006-01-02", req.AttendDate)
	if err != nil {
		return nil, errors.New("invalid attendance date format")
	}
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	if attendDate.After(today) {
		hasTransferredLesson := s.attendanceRepo.IsHaveTransferredLesson(req.GroupId)
		if hasTransferredLesson {
			cutoffTime := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
			if attendDate.Equal(today.AddDate(0, 0, -1)) && now.After(cutoffTime) {
				return nil, errors.New("attendance cannot be set for yesterday after 12 PM")
			}

			if req.Status == -1 {
				err = s.attendanceRepo.DeleteAttendance(req.GroupId, req.StudentId, req.TeacherId, req.AttendDate)
				if err != nil {
					return nil, err
				}
				return &pb.AbsResponse{
					Status:  200,
					Message: "Attendance successfully deleted",
				}, nil
			} else {
				err = s.attendanceRepo.CreateAttendance(ctx, companyId, req.GroupId, req.StudentId, req.TeacherId, req.AttendDate, req.Status, req.ActionById, req.ActionByRole)
				if err != nil {
					return nil, err
				}
				return &pb.AbsResponse{
					Status:  200,
					Message: "Attendance successfully created",
				}, nil
			}
		}
		return nil, errors.New("attendance date cannot be in the future")
	}
	validDay, err := s.attendanceRepo.IsValidGroupDay(ctx, req.GroupId, today)
	if err != nil {
		return nil, err
	}
	if !validDay {
		return nil, errors.New("attendance cannot be created today; group is not active")
	}

	if req.Status == -1 {
		err = s.attendanceRepo.DeleteAttendance(req.GroupId, req.StudentId, req.TeacherId, req.AttendDate)
		if err != nil {
			return nil, err
		}
		return &pb.AbsResponse{
			Status:  200,
			Message: "Attendance successfully deleted",
		}, nil
	} else {
		err = s.attendanceRepo.CreateAttendance(ctx, companyId, req.GroupId, req.StudentId, req.TeacherId, req.AttendDate, req.Status, req.ActionById, req.ActionByRole)
		if err != nil {
			return nil, err
		}
		return &pb.AbsResponse{
			Status:  200,
			Message: "Attendance successfully created",
		}, nil
	}
}
func (s *AttendanceService) CalculateTeacherSalaryByAttendance(ctx context.Context, req *pb.CalculateTeacherSalaryRequest) (*pb.CalculateTeacherSalaryResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}

	var response []*pb.AbsCalculateSalary
	groups := s.attendanceRepo.GetAllGroupsByTeacherId(req.TeacherId, req.From, req.To)
	for _, group := range groups {
		var absCalculate pb.AbsCalculateSalary
		absCalculate.GroupId = group.Id
		absCalculate.GroupName = group.Name
		absCalculate.CommonLessonCountInPeriod = group.LessonCountOnPeriod

		attendancesMap, err := s.attendanceRepo.GetAttendanceByTeacherAndGroup(companyId, req.TeacherId, group.Id, req.From, req.To)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error while getting attendance by teacher and group on calculating %v", err.Error()))
		}

		for studentId, attendances := range attendancesMap {
			var passedLessonCount int32
			var totalSalary float32
			var studentName string
			var priceType string
			var totalCount float64
			var coursePrice float64
			for _, attendance := range attendances {
				passedLessonCount++
				totalSalary += attendance.Price
				studentName = attendance.StudentName
				priceType = attendance.PriceType
				totalCount = attendance.TotalCount
				coursePrice = attendance.CoursePrice
			}
			absCalculate.Salaries = append(absCalculate.Salaries, &pb.StudentSalary{
				StudentId:                studentId,
				StudentName:              studentName,
				PassedLessonCount:        passedLessonCount,
				CalculatedSalaryInPeriod: int32(totalSalary),
				PriceType:                priceType,
				TotalCount:               totalCount,
				CoursePrice:              coursePrice,
			})
		}

		response = append(response, &absCalculate)
	}

	return &pb.CalculateTeacherSalaryResponse{Salaries: response}, nil
}
