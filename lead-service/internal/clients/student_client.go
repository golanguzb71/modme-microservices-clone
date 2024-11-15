package clients

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"lid-service/proto/pb"
)

type StudentClient struct {
	client pb.StudentServiceClient
}

func NewStudentClient(addr string) *StudentClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	fmt.Println(addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := pb.NewStudentServiceClient(conn)
	return &StudentClient{client: client}
}

func (gc *StudentClient) CreateStudent(ctx context.Context, phoneNumber, name, dateBirth, groupId, dateFrom, createdBy string, gender bool) (*pb.AbsResponse, error) {
	req := pb.CreateStudentRequest{
		PhoneNumber:       phoneNumber,
		Name:              name,
		DateOfBirth:       dateBirth,
		Gender:            gender,
		AdditionalContact: "",
		Address:           "",
		PassportId:        "",
		TelegramUsername:  "",
		GroupId:           groupId,
		DateFrom:          dateFrom,
		CreatedBy:         createdBy,
	}
	resp, err := gc.client.CreateStudent(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
