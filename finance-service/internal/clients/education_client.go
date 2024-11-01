package clients

import (
	"context"
	"finance-service/proto/pb"
	"google.golang.org/grpc"
)

type EducationClient struct {
	studentClient pb.StudentServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	studentClient := pb.NewStudentServiceClient(conn)
	return &EducationClient{studentClient: studentClient}, nil
}

func (ec *EducationClient) GetStudentById(studentId string) (string, string, error) {
	student, err := ec.studentClient.GetStudentById(context.TODO(), &pb.NoteStudentByAbsRequest{Id: studentId})
	if err != nil {
		return "", "", err
	}
	return student.Name, student.Phone, nil
}

func (ec *EducationClient) GetStudentsByGroupId(groupId string) (*pb.GetStudentsByGroupIdResponse, error) {
	return ec.studentClient.GetStudentsByGroupId(context.TODO(), &pb.GetStudentsByGroupIdRequest{GroupId: groupId, WithOutdated: true})
}

func (ec *EducationClient) ChangeUserBalanceHistory(studentId string, amount string, givenDate string, comment string, paymentType string, actionById string, actionByName string, groupId string) error {
	_, err := ec.studentClient.ChangeUserBalanceHistory(context.TODO(), &pb.ChangeUserBalanceHistoryRequest{
		StudentId:     studentId,
		Amount:        amount,
		GivenDate:     givenDate,
		Comment:       comment,
		PaymentType:   paymentType,
		CreatedBy:     actionById,
		CreatedByName: actionByName,
		GroupId:       groupId,
	})
	if err != nil {
		return err
	}
	return nil
}
