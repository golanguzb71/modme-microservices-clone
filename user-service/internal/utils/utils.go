package utils

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	fmt.Println(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println("company id checking ............")
	if ok {
		fmt.Println("company id ok is true")
		if _, ok := md["company_id"]; ok {
			fmt.Println("company id md company_id ok very true ")
			fmt.Println("here the value ", md["company_id"][0])
			return md["company_id"][0]
		}
	}
	fmt.Println("company id topilmadi")
	return ""
}
