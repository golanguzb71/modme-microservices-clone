package clients

import (
	"context"
	"education-service/proto/pb"
	"google.golang.org/grpc"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{client: client}, nil
}

func (uc *UserClient) GetTeacherById(ctx context.Context, teacherId string) (string, error) {
	user, err := uc.client.GetUserById(ctx, &pb.UserAbsRequest{UserId: teacherId})
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

func (uc *UserClient) GetUserByCompanyId(ctx context.Context, companyId, role string) (*pb.GetUserByCompanyIdResponse, error) {
	return uc.client.GetUserByCompanyId(ctx, &pb.GetUserByCompanyIdRequest{CompanyId: companyId, Role: role})
}
