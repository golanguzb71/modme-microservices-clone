package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func IsValidLessonDay(db *sql.DB, groupId, fromDate string) (bool, error) {
	var lessonDays []string

	err := db.QueryRow(`SELECT days FROM groups WHERE id = $1`, groupId).Scan(pq.Array(&lessonDays))
	if err != nil {
		return false, fmt.Errorf("failed to retrieve lesson days: %v", err)
	}
	parsedDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return false, fmt.Errorf("invalid date format for 'fromDate': %v", err)
	}

	dayOfWeek := parsedDate.Weekday().String()
	switch dayOfWeek {
	case "Monday":
		dayOfWeek = "DUSHANBA"
	case "Tuesday":
		dayOfWeek = "SESHANBA"
	case "Wednesday":
		dayOfWeek = "CHORSHANBA"
	case "Thursday":
		dayOfWeek = "PAYSHANBA"
	case "Friday":
		dayOfWeek = "JUMA"
	case "Saturday":
		dayOfWeek = "SHANBA"
	case "Sunday":
		dayOfWeek = "YAKSHANBA"
	}
	for _, lessonDay := range lessonDays {
		if strings.ToUpper(lessonDay) == strings.ToUpper(dayOfWeek) {
			return true, nil
		}
	}

	return false, nil
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
