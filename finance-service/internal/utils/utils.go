package utils

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
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

func ParseAmount(sum string) (float64, error) {
	amount, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		return 0, errors.New("invalid amount format")
	}
	return amount, nil
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

func NewTimoutContext(ctx context.Context, companyId string) (context.Context, context.CancelFunc) {
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, cancelFunc := context.WithTimeout(ctx, time.Second*60)
	return res, cancelFunc
}
