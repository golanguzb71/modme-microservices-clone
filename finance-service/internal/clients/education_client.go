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
	return ec.studentClient.GetStudentsByGroupId(context.TODO(), &pb.GetStudentsByGroupIdRequest{GroupId: groupId, WithOutdated: false})
}
