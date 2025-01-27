package repository

import (
	"context"
	"database/sql"
	"errors"
	"finance-service/internal/clients"
	"finance-service/internal/utils"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type TeacherSalaryRepository struct {
	db         *sql.DB
	userClient *clients.UserClient
}

func (r *TeacherSalaryRepository) CreateTeacherSalary(ctx context.Context, companyId string, amount int32, teacherId string, amountType string) (*pb.AbsResponse, error) {
	if amountType == "PERCENT" && (amount > 100 || amount < 0) {
		return nil, status.Errorf(codes.Aborted, "invalid amount for PERCENT: must be between 0 and 100")
	}
	if amountType != "PERCENT" && amount < 10000 {
		return nil, status.Errorf(codes.Aborted, "invalid amount: must be non-negative")
	}
	_, err := r.db.Exec("INSERT INTO teacher_salary (teacher_id, salary_type, salary_type_count , company_id) VALUES ($1, $2, $3 , $4)", teacherId, amountType, amount, companyId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert data: %v", err)
	}
	return &pb.AbsResponse{
		Status:  http.StatusCreated,
		Message: "created",
	}, nil
}

func (r *TeacherSalaryRepository) DeleteTeacherSalary(ctx context.Context, companyId string, teacherId string) (*pb.AbsResponse, error) {
	result, err := r.db.Exec("DELETE FROM teacher_salary WHERE teacher_id = $1 and company_id=$2", teacherId, companyId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete salary: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "no salary found with the provided teacherId")
	}

	return &pb.AbsResponse{Message: "Salary deleted successfully"}, nil
}

func (r *TeacherSalaryRepository) GetTeacherSalary(ctx context.Context, companyId string) (*pb.GetTeachersSalaryRequest, error) {
	rows, err := r.db.Query("SELECT teacher_id, salary_type, salary_type_count  FROM teacher_salary where company_id=$1", companyId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve salaries: %v", err)
	}
	defer rows.Close()
	var salaries []*pb.AbsGetTeachersSalary
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for rows.Next() {
		var teacherId, salaryType, teacherName string
		var amount int32

		if err := rows.Scan(&teacherId, &salaryType, &amount); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan row: %v", err)
		}
		user, err := r.userClient.GetUserById(ctx, teacherId)
		if err != nil {
			teacherName = "Teacher Name not available"
		}
		teacherName = user.Name
		salaries = append(salaries, &pb.AbsGetTeachersSalary{
			TeacherId:   teacherId,
			Type:        salaryType,
			Amount:      amount,
			TeacherName: teacherName,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error iterating rows: %v", err)
	}

	return &pb.GetTeachersSalaryRequest{Salaries: salaries}, nil
}

func (r *TeacherSalaryRepository) GetTeacherSalaryByTeacherID(ctx context.Context, companyId string, teacherId string) (*pb.AbsGetTeachersSalary, error) {
	var salary pb.AbsGetTeachersSalary
	err := r.db.QueryRow("SELECT teacher_id, salary_type, salary_type_count FROM teacher_salary WHERE teacher_id = $1 and company_id=$2", teacherId, companyId).
		Scan(&salary.TeacherId, &salary.Type, &salary.Amount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "salary not found for teacherId: %s", teacherId)
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve salary: %v", err)
	}

	return &salary, nil
}

func NewTeacherSalaryRepository(db *sql.DB, userClient *clients.UserClient) *TeacherSalaryRepository {
	return &TeacherSalaryRepository{db: db, userClient: userClient}
}
