package utils

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

func EncodePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func RecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in gRPC call: %v\n", r)
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()
	return handler(ctx, req)
}

func GetCompanyDetails(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if _, ok := md["company_id"]; ok {
			fmt.Println("here the value ", md["company_id"][0])
			return md["company_id"][0]
		}
	}
	return ""
}

func NewTimoutContext(ctx context.Context, companyId string) (context.Context, context.CancelFunc) {
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, cancelFunc := context.WithTimeout(ctx, time.Second*15)
	return res, cancelFunc
}
