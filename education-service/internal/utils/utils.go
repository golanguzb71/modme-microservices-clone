package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"math"
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

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in gRPC call: %v\n", r)
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()
	return handler(ctx, req)
}

func CalculateMoneyForStatus(db *sql.DB, manualPriceForCourse *float64, groupId string, tillDate string) (float64, error) {
	var coursePrice float64
	var courseDurationLesson int
	var groupStartTime string
	var groupDays []string
	var dateType string

	query := `
        SELECT c.price, c.course_duration, g.start_time, g.days, g.date_type
        FROM groups g
        JOIN courses c ON g.course_id = c.id
        WHERE g.id = $1
    `

	err := db.QueryRow(query, groupId).Scan(&coursePrice, &courseDurationLesson, &groupStartTime, pq.Array(&groupDays), &dateType)
	if err != nil {
		return 0, fmt.Errorf("error getting course details: %v", err)
	}

	if manualPriceForCourse != nil {
		coursePrice = coursePrice - *manualPriceForCourse
	}

	tillDateParsed, err := time.Parse("2006-01-02", tillDate)
	if err != nil {
		return 0, fmt.Errorf("error parsing till date: %v", err)
	}

	endOfMonth := time.Date(tillDateParsed.Year(), tillDateParsed.Month(), 1, 23, 59, 59, 999999999, tillDateParsed.Location()).AddDate(0, 1, -1)

	totalLessonsInMonth := calculateLessonsInMonth(groupDays, dateType, time.Date(tillDateParsed.Year(), tillDateParsed.Month(), 1, 0, 0, 0, 0, tillDateParsed.Location()), endOfMonth)
	if totalLessonsInMonth == 0 {
		return 0, fmt.Errorf("no lessons scheduled for the given month, avoiding division by zero")
	}

	remainingLessons := calculateRemainingLessons(groupDays, dateType, tillDateParsed, endOfMonth)
	if remainingLessons > totalLessonsInMonth {
		remainingLessons = totalLessonsInMonth
	}

	remainingMoney := coursePrice / float64(totalLessonsInMonth) * float64(remainingLessons)

	if remainingMoney < 0 {
		remainingMoney = math.Ceil(remainingMoney)
	} else {
		remainingMoney = math.Floor(remainingMoney)
	}

	return remainingMoney, nil
}
func calculateLessonsInMonth(groupDays []string, dateType string, startDate, endDate time.Time) int {
	totalLessons := 0
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if isLessonDay(currentDate, groupDays, dateType) {
			totalLessons++
		}
	}
	return totalLessons
}

func calculateRemainingLessons(groupDays []string, dateType string, currentDate, endDate time.Time) int {
	remainingLessons := 0
	for ; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if isLessonDay(currentDate, groupDays, dateType) {
			remainingLessons++
		}
	}
	return remainingLessons
}

func isLessonDay(currentDate time.Time, groupDays []string, dateType string) bool {
	dayName := getDayName(currentDate.Weekday())
	for _, groupDay := range groupDays {
		if groupDay == dayName {
			switch dateType {
			case "JUFT":
				return currentDate.Day()%2 == 0
			case "TOQ":
				return currentDate.Day()%2 != 0
			default:
				return true
			}
		}
	}
	return false
}

func getDayName(weekday time.Weekday) string {
	days := map[time.Weekday]string{
		time.Monday:    "DUSHANBA",
		time.Tuesday:   "SESHANBA",
		time.Wednesday: "CHORSHANBA",
		time.Thursday:  "PAYSHANBA",
		time.Friday:    "JUMA",
		time.Saturday:  "SHANBA",
		time.Sunday:    "YAKSHANBA",
	}
	return days[weekday]
}

func CheckGroupAndTeacher(db *sql.DB, groupId, actionRole string, actionId string) bool {
	if actionRole == "TEACHER" {
		fmt.Println(actionRole)
		checker := false
		err := db.QueryRow(`SELECT exists(SELECT 1 FROM groups where id=$1 and teacher_id=$2)`, groupId, actionId).Scan(&checker)
		if err != nil || !checker {
			return false
		}
	} else if actionRole == "EMPLOYEE" {
		return false
	}
	return true
}

func CalculateMoneyForLesson(db *sql.DB, price *float64, studentId string, groupId string, attendDate string, discountAmount *float64) error {
	var coursePrice float64
	err := db.QueryRow(`SELECT price FROM courses c join groups g on c.id=g.course_id where g.id=$1`, groupId).Scan(&coursePrice)
	if err != nil {
		return err
	}
	if discountAmount != nil {
		coursePrice = coursePrice - *discountAmount
	}
	parsedDate, err := time.Parse("2006-01-02", attendDate)
	if err != nil {
		return err
	}
	firstOfMonth := time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, parsedDate.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	var lessonCount int
	err = db.QueryRow(`
        WITH RECURSIVE dates AS (
            SELECT $1::date AS date
            UNION ALL
            SELECT date + 1
            FROM dates
            WHERE date < $2
        ),
        valid_days AS (
            SELECT d.date
            FROM dates d
            CROSS JOIN (
                SELECT unnest(days) AS day_name,
                       start_date,
                       end_date
                FROM groups
                WHERE id = $3
            ) g
            WHERE 
                CASE 
                    WHEN EXTRACT(DOW FROM d.date) = 1 THEN 'DUSHANBA'
                    WHEN EXTRACT(DOW FROM d.date) = 2 THEN 'SESHANBA'
                    WHEN EXTRACT(DOW FROM d.date) = 3 THEN 'CHORSHANBA'
                    WHEN EXTRACT(DOW FROM d.date) = 4 THEN 'PAYSHANBA'
                    WHEN EXTRACT(DOW FROM d.date) = 5 THEN 'JUMA'
                    WHEN EXTRACT(DOW FROM d.date) = 6 THEN 'SHANBA'
                    WHEN EXTRACT(DOW FROM d.date) = 0 THEN 'YAKSHANBA'
                END = g.day_name
                AND d.date >= g.start_date
                AND d.date <= g.end_date
        )
        SELECT COUNT(*) 
        FROM valid_days
    `, firstOfMonth, lastOfMonth, groupId).Scan(&lessonCount)
	if err != nil {
		return err
	}

	if lessonCount == 0 {
		return fmt.Errorf("no lessons found in the month")
	}

	*price = math.Round(coursePrice / float64(lessonCount))
	return nil
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
