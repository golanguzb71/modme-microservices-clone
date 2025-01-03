package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"user-service/internal/clients"
	"user-service/internal/utils"
	"user-service/proto/pb"
)

type UserRepository struct {
	db              *sql.DB
	groupClient     *clients.GroupClient
	groupClientChan chan *clients.GroupClient
}

func NewUserRepository(db *sql.DB, clientChan chan *clients.GroupClient) *UserRepository {
	return &UserRepository{db: db, groupClientChan: clientChan}
}

func (r *UserRepository) ensureGroupClient() error {
	if r.groupClient == nil {
		select {
		case client := <-r.groupClientChan:
			r.groupClient = client
		case <-time.After(5 * time.Second):
			return fmt.Errorf("failed to initialize GroupClient within timeout")
		}
	}
	return nil
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
	if err := r.ensureGroupClient(); err != nil {
		return nil, err
	}

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
	if role != "TEACHER" && role != "ADMIN" && role != "EMPLOYEE" && role != "CEO" {
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
	var role string
	err := r.db.QueryRow(`SELECT role FROM users where id=$1`, id).Scan(&role)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}
	if role == "TEACHER" {
		groupCount, err := r.groupClient.GetGroupsByTeacherId(id, false)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		} else if groupCount > 0 {
			return nil, status.Errorf(codes.DataLoss, "Ushbu teacherga bog'langan active guruhlar mavjud iltimos avval guruhni arxivlang !!")
		}
	} else if role == "CEO" {
		return nil, status.Errorf(codes.Canceled, "Tizimda CEO bo'lishi shart !!")
	}
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
		return nil, status.Errorf(codes.NotFound, "User not found")
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
func (r *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*pb.GetUserByIdResponse, string, error) {
	res := pb.GetUserByIdResponse{}
	var password string
	err := r.db.QueryRow(`SELECT id,
       full_name,
       phone_number,
       password,
       role,
       birth_date,
       gender,
       is_deleted,
       created_at FROM users where phone_number=$1`, phoneNumber).Scan(&res.Id, &res.Name, &res.PhoneNumber, &password, &res.Role, &res.BirthDate, &res.Gender, &res.IsDeleted, &res.CreatedAt)
	if err != nil {
		return nil, "", err
	}
	return &res, password, nil
}
func (r *UserRepository) GetAllStuff(isArchived bool) (*pb.GetAllStuffResponse, error) {
	query := `SELECT id, phone_number, role, full_name, birth_date, gender, is_deleted, created_at
              FROM users WHERE is_deleted = $1`
	rows, err := r.db.Query(query, isArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response pb.GetAllStuffResponse
	for rows.Next() {
		var user pb.GetUserByIdResponse
		err = rows.Scan(&user.Id, &user.PhoneNumber, &user.Role, &user.Name, &user.BirthDate, &user.Gender, &user.IsDeleted, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		response.Stuff = append(response.Stuff, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &response, nil
}
func (r *UserRepository) GetHistoryByUserId(id string) (*pb.GetHistoryByUserIdResponse, error) {
	query := `
	SELECT 
		updated_field,
		old_value,
		current_value,
		created_at
	FROM users_history
	WHERE user_id = $1
	ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user history: %v", err)
	}
	defer rows.Close()

	var historyItems []*pb.AbsGetHistoryByUserIdResponse
	for rows.Next() {
		var item pb.AbsGetHistoryByUserIdResponse
		if err := rows.Scan(&item.UpdatedField, &item.OldValue, &item.CurrentValue, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning history item: %v", err)
		}
		historyItems = append(historyItems, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through history rows: %v", err)
	}

	return &pb.GetHistoryByUserIdResponse{
		Histories: historyItems,
	}, nil
}

func (r *UserRepository) UpdateUserPassword(userId string, password string) (*pb.AbsResponse, error) {
	var userExists bool
	err := r.db.QueryRow(`SELECT exists(SELECT 1 FROM users where id=$1)`, userId).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, status.Errorf(codes.AlreadyExists, "user not found")
	}
	newEncodedPass, err := utils.EncodePassword(password)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Exec(`UPDATE users set password=$1 where id=$2`, newEncodedPass, userId)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "password updated",
	}, nil
}
