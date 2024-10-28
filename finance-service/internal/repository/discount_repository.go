package repository

import (
	"database/sql"
	"finance-service/internal/clients"
	"finance-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscountRepository struct {
	db            *sql.DB
	studentClient *clients.EducationClient
}

func (r *DiscountRepository) CreateDiscount(groupId string, studentId string, discountPrice, comment string) error {
	var checker bool
	_ = r.db.QueryRow(`SELECT exists(SELECT 1 FROM student_discount where group_id=$1 and student_id=$1)`, groupId, studentId).Scan(&checker)
	if checker {
		return status.Errorf(codes.AlreadyExists, "discount already exits")
	}
	_, err := r.db.Exec(`INSERT INTO student_discount(student_id, discount, group_id, comment) values ($1 ,$2 ,$3 , $4)`, studentId, discountPrice, groupId, comment)
	if err != nil {
		return status.Errorf(codes.Internal, "Status discount insert error %f", err)
	}
	return nil
}

func (r *DiscountRepository) DeleteDiscount(groupId string, studentId string) error {
	_, err := r.db.Exec(`DELETE FROM student_discount where group_id=$1 and student_id=$2`, groupId, studentId)
	if err != nil {
		return status.Errorf(codes.Internal, "Error while deleting disount %f", err)
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
		r.db.QueryRow(`SELECT  discount, comment, created_at FROM student_discount WHERE group_id=$1 and student_id=$2`, groupId, el.Id).Scan(&element.Discount, &element.Cause, &element.CreatedAt)
		res = append(res, &element)
	}
	result.Discounts = res
	return &result, nil
}

func NewDiscountRepository(db *sql.DB, studentClient *clients.EducationClient) *DiscountRepository {
	return &DiscountRepository{db: db, studentClient: studentClient}
}
