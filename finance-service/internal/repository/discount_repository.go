package repository

import (
	"database/sql"
	"finance-service/internal/clients"
	"finance-service/proto/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscountRepository struct {
	db            *sql.DB
	studentClient *clients.EducationClient
}

func (r *DiscountRepository) CreateDiscount(groupId string, studentId string, discountPrice, comment, startDate, endDate string, withTeacher bool) error {
	var checker bool
	err := r.db.QueryRow(`SELECT exists(SELECT 1 FROM student_discount WHERE group_id=$1 AND student_id=$2)`, groupId, studentId).Scan(&checker)
	if err != nil {
		return status.Errorf(codes.Internal, "Error checking existing discount: %v", err)
	}

	if checker {
		return status.Errorf(codes.AlreadyExists, "Discount already exists")
	}
	_, err = r.db.Exec(`INSERT INTO student_discount (student_id, discount, group_id, comment, start_at, end_at, withteacher) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, studentId, discountPrice, groupId, comment, startDate, endDate, withTeacher)
	if err != nil {
		return status.Errorf(codes.Internal, "Status discount insert error: %v", err)
	}
	_, err = r.db.Exec(`INSERT INTO student_discount_history (id, student_id, group_id, start_at, end_at, withteacher, comment, action)
		VALUES ($1, $2, $3, $4, $5 , $6 , $7, $8)`, uuid.New(), studentId, groupId, startDate, endDate, withTeacher, comment, "CREATE")
	if err != nil {
		return status.Errorf(codes.Internal, "Error inserting into student history: %v", err)
	}

	return nil
}
func (r *DiscountRepository) DeleteDiscount(groupId string, studentId string) error {
	var discount pb.AbsDiscountRequest
	var createdAt string
	err := r.db.QueryRow(`SELECT student_id, discount, group_id, comment, start_at, end_at, withteacher, created_at FROM student_discount where student_id=$1 and group_id=$2`, studentId, groupId).Scan(
		&discount.StudentId,
		&discount.DiscountPrice,
		&discount.GroupId,
		&discount.Comment,
		&discount.StartDate,
		&discount.EndDate,
		&discount.WithTeacher,
		&createdAt,
	)
	if err != nil {
		return status.Errorf(codes.NotFound, err.Error())
	}
	_, err = r.db.Exec(`DELETE FROM student_discount where group_id=$1 and student_id=$2`, groupId, studentId)
	if err != nil {
		return status.Errorf(codes.Internal, "Error while deleting disount %f", err)
	}
	_, err = r.db.Exec(`INSERT INTO student_discount_history (id, student_id, group_id, start_at, end_at, withteacher, comment, action)
		VALUES ($1, $2, $3, $4, $5 , $6 , $7, $8)`, uuid.New(), studentId, groupId, discount.StartDate, discount.EndDate, discount.WithTeacher, discount.Comment, "DELETE")
	if err != nil {
		return status.Errorf(codes.Internal, "Error inserting into student history: %v", err)
	}
	return nil
}
func (r *DiscountRepository) GetAllDiscountByGroup(groupId string) (*pb.GetInformationDiscountResponse, error) {
	resp, err := r.studentClient.GetStudentsByGroupId(groupId)
	if err != nil {
		return nil, err
	}
	result := pb.GetInformationDiscountResponse{}
	var res []*pb.AbsStudentDiscount
	students := resp.Students
	for _, el := range students {
		var element pb.AbsStudentDiscount
		element.StudentId = el.Id
		element.StudentName = el.Name
		element.StudentPhoneNumber = el.PhoneNumber
		r.db.QueryRow(`SELECT  discount, comment, created_at , start_at , end_at , withteacher FROM student_discount WHERE group_id=$1 and student_id=$2`, groupId, el.Id).Scan(&element.Discount, &element.Cause, &element.CreatedAt, &element.StartAt, &element.EndAt, &element.WithTeacher)
		res = append(res, &element)
	}
	result.Discounts = res
	return &result, nil
}
func NewDiscountRepository(db *sql.DB, studentClient *clients.EducationClient) *DiscountRepository {
	return &DiscountRepository{db: db, studentClient: studentClient}
}
