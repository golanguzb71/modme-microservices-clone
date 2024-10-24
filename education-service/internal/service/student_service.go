package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/proto/pb"
)

type StudentService struct {
	pb.UnimplementedStudentServiceServer
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetAllStudent(ctx context.Context, req *pb.GetAllStudentRequest) (*pb.GetAllStudentResponse, error) {
	return s.repo.GetAllStudent(req.Condition, req.Page, req.Size)
}

func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.AbsResponse, error) {
	if err := s.repo.CreateStudent(req.CreatedBy, req.PhoneNumber, req.Name, req.GroupId, req.Address, req.AdditionalContact, req.DateFrom, req.DateOfBirth, req.Gender, req.PassportId, req.TelegramUsername); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "student created successfully",
	}, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.AbsResponse, error) {
	if err := s.repo.UpdateStudent(req.StudentId, req.PhoneNumber, req.Name, req.Address, req.AdditionalContact, req.DateOfBirth, req.Gender, req.PassportId); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "student updated successfully",
	}, nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	if err := s.repo.DeleteStudent(req.Id); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "accomplished",
	}, nil
}

func (s *StudentService) AddToGroup(ctx context.Context, req *pb.AddToGroupRequest) (*pb.AbsResponse, error) {
	if err := s.repo.AddToGroup(req.GroupId, req.StudentIds, req.CreatedDate, req.CreatedBy); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "students added to group",
	}, nil
}

func (s *StudentService) GetStudentById(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetStudentByIdResponse, error) {
	return s.repo.GetStudentById(req.Id)
}
func (s *StudentService) GetNoteByStudent(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetNotesByStudent, error) {
	return s.repo.GetNoteByStudent(req.Id)
}
func (s *StudentService) CreateNoteForStudent(ctx context.Context, req *pb.CreateNoteRequest) (*pb.AbsResponse, error) {
	return s.repo.CreateNoteForStudent(req.Note, req.StudentId)
}
func (s *StudentService) DeleteStudentNote(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.AbsResponse, error) {
	return s.repo.DeleteStudentNote(req.Id)
}
func (s *StudentService) SearchStudent(ctx context.Context, req *pb.SearchStudentRequest) (*pb.SearchStudentResponse, error) {
	return s.repo.SearchStudent(req.Value)
}
func (s *StudentService) GetHistoryGroupById(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetHistoryGroupResponse, error) {
	return s.repo.GetHistoryGroupById(req.Id)
}
func (s *StudentService) GetHistoryStudentById(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetHistoryStudentResponse, error) {
	return s.repo.GetHistoryStudentById(req.Id)
}
func (s *StudentService) TransferLessonDate(ctx context.Context, req *pb.TransferLessonRequest) (*pb.AbsResponse, error) {
	return s.repo.TransferLessonDate(req.GroupId, req.From, req.To)
}

func (s *StudentService) ChangeConditionStudent(ctx context.Context, req *pb.ChangeConditionStudentRequest) (*pb.AbsResponse, error) {
	return s.repo.ChangeConditionStudent(req.StudentId, req.GroupId, req.Status, req.ReturnTheMoney, req.TillDate)
}
