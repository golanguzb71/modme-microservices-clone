package repository

import (
	"context"
	"database/sql"
	"education-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetAllStudent(condition string, page string, size string) (*pb.GetAllStudentResponse, error) {
	return nil, nil
}

func (r *StudentRepository) CreateStudent(ctx context.Context, phoneNumber string, name string, groupId string, address string, additionalContact string, dateFrom string, birthDate string, gender bool, passportId string, telegramUsername string) error {
	createdBy, ok := ctx.Value("createdBy").(string)
	if !ok {
		return fmt.Errorf("could not retrieve createdBy from context or invalid type")
	}
	studentId := uuid.New()
	_, err := r.db.Exec(`INSERT INTO students(id, name, phone, date_of_birth, gender, telegram_username, passport_id, additional_contact, address) values ($1, $2,$3,$4,$5,$6,$7,$8,$9)`, studentId, name, phoneNumber, birthDate, gender, telegramUsername, passportId, additionalContact, address)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`INSERT INTO group_students(id, group_id, student_id, created_by) values ($1 ,$2 ,$3 ,$4)`, uuid.New(), groupId, studentId, createdBy)
	if err != nil {
		return err
	}
	return nil
}

func (r *StudentRepository) UpdateStudent(studentId string, number string, name string, address string, additionalContact string, birth string, gender bool, passportId string) error {
	return nil
}

func (r *StudentRepository) DeleteStudent(studentId string) error {
	return nil
}

func (r *StudentRepository) AddToGroup(groupId string, studentIds []string, createdDate string) error {
	return nil
}
