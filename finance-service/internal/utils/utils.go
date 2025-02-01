package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
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

const (
	botToken = "7667139311:AAECrJwO0cYWnIx8AmNuH8O_hZezVUufsuI"
	chatID   = "6805374430"
)

func SendTelegramMessage(message string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"
	data := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return err
}
