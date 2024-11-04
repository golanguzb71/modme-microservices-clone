package repository

import (
	"database/sql"
	"errors"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type TeacherSalaryRepository struct {
	db *sql.DB
}

func (r *TeacherSalaryRepository) CreateTeacherSalary(amount int32, teacherId string, amountType string) (*pb.AbsResponse, error) {
	if amountType == "PERCENT" && (amount > 100 || amount < 0) {
		return nil, status.Errorf(codes.Aborted, "invalid amount for PERCENT: must be between 0 and 100")
	}
	if amountType != "PERCENT" && amount < 10000 {
		return nil, status.Errorf(codes.Aborted, "invalid amount: must be non-negative")
	}
	_, err := r.db.Exec("INSERT INTO teacher_salary (teacher_id, salary_type, salary_type_count) VALUES ($1, $2, $3)", teacherId, amountType, amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert data: %v", err)
	}
	return &pb.AbsResponse{
		Status:  http.StatusCreated,
		Message: "created",
	}, nil
}

func (r *TeacherSalaryRepository) DeleteTeacherSalary(teacherId string) (*pb.AbsResponse, error) {
	result, err := r.db.Exec("DELETE FROM teacher_salary WHERE teacher_id = $1", teacherId)
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

func (r *TeacherSalaryRepository) GetTeacherSalary() (*pb.GetTeachersSalaryRequest, error) {
	rows, err := r.db.Query("SELECT teacher_id, salary_type, salary_type_count FROM teacher_salary")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve salaries: %v", err)
	}
	defer rows.Close()

	var salaries []*pb.AbsGetTeachersSalary
	for rows.Next() {
		var teacherId, salaryType string
		var amount int32

		if err := rows.Scan(&teacherId, &salaryType, &amount); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan row: %v", err)
		}

		salaries = append(salaries, &pb.AbsGetTeachersSalary{
			TeacherId: teacherId,
			Type:      salaryType,
			Amount:    amount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error iterating rows: %v", err)
	}

	return &pb.GetTeachersSalaryRequest{Salaries: salaries}, nil
}

func (r *TeacherSalaryRepository) GetTeacherSalaryByTeacherID(teacherId string) (*pb.AbsGetTeachersSalary, error) {
	var salary pb.AbsGetTeachersSalary
	err := r.db.QueryRow("SELECT teacher_id, salary_type, salary_type_count FROM teacher_salary WHERE teacher_id = $1", teacherId).
		Scan(&salary.TeacherId, &salary.Type, &salary.Amount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "salary not found for teacherId: %s", teacherId)
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve salary: %v", err)
	}

	return &salary, nil
}

func NewTeacherSalaryRepository(db *sql.DB) *TeacherSalaryRepository {
	return &TeacherSalaryRepository{db: db}
}
