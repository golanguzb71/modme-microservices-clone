package clients

import (
	"context"
	"finance-service/proto/pb"
	"google.golang.org/grpc"
)

type EducationClient struct {
	studentClient pb.StudentServiceClient
	groupClient   pb.GroupServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	studentClient := pb.NewStudentServiceClient(conn)
	groupClient := pb.NewGroupServiceClient(conn)
	return &EducationClient{studentClient: studentClient, groupClient: groupClient}, nil
}

func (ec *EducationClient) GetStudentById(studentId string) (string, string, float64, error) {
	student, err := ec.studentClient.GetStudentById(context.TODO(), &pb.NoteStudentByAbsRequest{Id: studentId})
	if err != nil {
		return "", "", 0, err
	}
	return student.Name, student.Phone, student.Balance, nil
}

func (ec *EducationClient) GetStudentsByGroupId(groupId string) (*pb.GetStudentsByGroupIdResponse, error) {
	return ec.studentClient.GetStudentsByGroupId(context.TODO(), &pb.GetStudentsByGroupIdRequest{GroupId: groupId, WithOutdated: true})
}

func (ec *EducationClient) GetGroupNameById(groupId string) string {
	group, err := ec.groupClient.GetGroupById(context.TODO(), &pb.GetGroupByIdRequest{Id: groupId})
	if err != nil {
		return ""
	}
	return group.Name
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

func (ec *EducationClient) ChangeUserBalanceHistoryByDebit(ctx context.Context, studentId, oldDebit, givenDate, comment, paymentType, createdBy, createdByName, groupId, currentDebit string) (*pb.AbsResponse, error) {
	return ec.studentClient.ChangeUserBalanceHistoryByDebit(ctx, &pb.ChangeUserBalanceHistoryByDebitRequest{
		StudentId:     studentId,
		OldDebit:      oldDebit,
		GivenDate:     givenDate,
		Comment:       comment,
		PaymentType:   paymentType,
		CreatedBy:     createdBy,
		CreatedByName: createdByName,
		GroupId:       groupId,
		CurrentDebit:  currentDebit,
	})
}

func (ec *EducationClient) GetGroupsAndCommentsByStudentId(ctx context.Context, studentId string) (*pb.GetGroupsByStudentResponse, error) {
	return ec.groupClient.GetGroupsByStudentId(ctx, &pb.StudentIdRequest{StudentId: studentId})
}
