package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
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
