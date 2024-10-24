package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
	"user-service/internal/clients"
	"user-service/internal/utils"
	"user-service/proto/pb"
)

type UserRepository struct {
	db          *sql.DB
	groupClient *clients.GroupClient
}

func NewUserRepository(db *sql.DB, client *clients.GroupClient) *UserRepository {
	return &UserRepository{db: db, groupClient: client}
}

func (r *UserRepository) CreateUser(gender bool, number string, birthDate string, name string, password string, role string) (*pb.AbsResponse, error) {
	encodedPassword, err := utils.EncodePassword(password)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(`INSERT INTO users(id, full_name, phone_number, password, role, birth_date, gender) values ($1 , $2 , $3 , $4 , $5 , $6, $7)`, uuid.New(), name, number, encodedPassword, role, birthDate, gender)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "created",
	}, nil
}
func (r *UserRepository) GetTeachers(isDeleted bool) (*pb.GetTeachersResponse, error) {
	rows, err := r.db.Query(`SELECT id, full_name, phone_number FROM users WHERE is_deleted=$1 AND role='TEACHER'`, isDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response pb.GetTeachersResponse
	for rows.Next() {
		var id, fullName, phoneNumber string
		if err := rows.Scan(&id, &fullName, &phoneNumber); err != nil {
			return nil, err
		}
		activeGroupsCount, err := r.groupClient.GetGroupsByTeacherId(id, false)
		if err != nil {
			return nil, err
		}
		response.Teachers = append(response.Teachers, &pb.AbsTeacher{
			Id:           id,
			FullName:     fullName,
			PhoneNumber:  phoneNumber,
			ActiveGroups: fmt.Sprintf("%d", activeGroupsCount),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *UserRepository) GetUserById(userId string) (*pb.GetUserByIdResponse, error) {
	var response pb.GetUserByIdResponse
	err := r.db.QueryRow(`SELECT id,
       full_name,
       phone_number,
       role,
       birth_date,
       gender,
       is_deleted,
       created_at FROM users where id=$1`, userId).Scan(&response.Id, &response.Name, &response.PhoneNumber, &response.Role, &response.BirthDate, &response.Gender, &response.IsDeleted, &response.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
func (r *UserRepository) UpdateUser(userId string, name string, gender bool, role string, birthDate string, phoneNumber string) (*pb.AbsResponse, error) {
	if role != "TEACHER" && role != "ADMIN" && role != "EMPLOYEE" {
		return &pb.AbsResponse{Status: 400, Message: "Invalid role"}, nil
	}
	query := `
        UPDATE users 
        SET full_name = $1, phone_number = $2, gender = $3, role = $4, birth_date = $5
        WHERE id = $6
    `

	_, err := r.db.Exec(query, name, phoneNumber, gender, role, birthDate, userId)
	if err != nil {
		return nil, err
	}

	return &pb.AbsResponse{Status: 200, Message: "User updated successfully"}, nil
}
func (r *UserRepository) DeleteUser(id string) (*pb.AbsResponse, error) {
	query := `
        UPDATE users 
        SET is_deleted = NOT is_deleted 
        WHERE id = $1
    `
	result, err := r.db.Exec(query, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return &pb.AbsResponse{Status: 404, Message: "User not found"}, nil
	}
	return &pb.AbsResponse{Status: 200, Message: "User status toggled successfully"}, nil
}
func (r *UserRepository) GetAllEmployee(isArchived bool) (*pb.GetAllEmployeeResponse, error) {
	query := `
        SELECT id, full_name, phone_number, role, birth_date, gender, is_deleted, created_at 
        FROM users 
        WHERE is_deleted = $1 AND role IN ('ADMIN', 'EMPLOYEE')
    `

	rows, err := r.db.Query(query, isArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*pb.GetUserByIdResponse

	for rows.Next() {
		var emp pb.GetUserByIdResponse
		var createdAt time.Time
		err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Role, &emp.BirthDate, &emp.Gender, &emp.IsDeleted, &createdAt)
		if err != nil {
			return nil, err
		}
		emp.CreatedAt = createdAt.Format(time.RFC3339)
		employees = append(employees, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllEmployeeResponse{Employees: employees}, nil
}
