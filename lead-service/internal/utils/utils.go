package utils

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

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
func GetCompanyId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if _, ok := md["company_id"]; ok {
			return md["company_id"][0]
		}
	}
	return ""
}

func NewTimoutContext(companyId string) (context.Context, context.CancelFunc) {
	var ctx context.Context
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, cancel := context.WithTimeout(ctx, time.Second*15)
	return res, cancel
}
