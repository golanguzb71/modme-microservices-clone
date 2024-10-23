package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"user-service/internal/utils"
	"user-service/proto/pb"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
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
	rows, err := r.db.Query(`SELECT id, full_name, phone_number FROM users where is_deleted=$1 and role='TEACHER'`, isDeleted)
	if err != nil {
		return nil, err
	}

}
