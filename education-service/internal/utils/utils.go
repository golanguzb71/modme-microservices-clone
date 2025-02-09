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
	var groupStartDate string
	var groupDays []string
	var dateType string

	query := `
        SELECT c.price, c.course_duration, g.start_date, g.days, g.date_type
        FROM groups g
        JOIN courses c ON g.course_id = c.id
        WHERE g.id = $1
    `

	err := db.QueryRow(query, groupId).Scan(&coursePrice, &courseDurationLesson, &groupStartDate, pq.Array(&groupDays), &dateType)
	if err != nil {
		return 0, fmt.Errorf("error getting course details: %v", err)
	}

	if manualPriceForCourse != nil {
		coursePrice = *manualPriceForCourse
	}

	tillDateParsed, err := time.Parse("2006-01-02", tillDate)
	if err != nil {
		return 0, fmt.Errorf("error parsing till date: %v", err)
	}
	// Get start and end of month
	startOfMonth := time.Date(tillDateParsed.Year(), tillDateParsed.Month(), 1, 0, 0, 0, 0, tillDateParsed.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	// Get all lesson dates for the month
	lessonDates := getLessonDatesInMonth(groupDays, dateType, startOfMonth, endOfMonth)
	if len(lessonDates) == 0 {
		return 0, fmt.Errorf("no lessons scheduled for the current month")
	}
	fmt.Println("lesson dates", lessonDates)
	fmt.Println("lesson dates lenghth", len(lessonDates))

	// Find first lesson date of the month
	firstLessonDate := lessonDates[0]

	// If we're before or on the first lesson date, return full amount
	if !tillDateParsed.After(firstLessonDate) {
		return coursePrice, nil
	}

	// Count passed lessons
	passedLessons := 0
	for _, lessonDate := range lessonDates {
		if lessonDate.Before(tillDateParsed) {
			passedLessons++
		}
	}

	// Calculate money per lesson and remaining amount
	pricePerLesson := coursePrice / float64(len(lessonDates))
	remainingMoney := coursePrice - (float64(passedLessons) * pricePerLesson)
	return math.Round(remainingMoney), nil
}

func getLessonDatesInMonth(groupDays []string, dateType string, startDate, endDate time.Time) []time.Time {
	var lessonDates []time.Time
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if isLessonDay(currentDate, groupDays, dateType) {
			lessonDates = append(lessonDates, currentDate)
		}
	}
	return lessonDates
}

func isLessonDay(currentDate time.Time, groupDays []string, dateType string) bool {
	dayName := getDayName(currentDate.Weekday())
	isGroupDay := false
	for _, groupDay := range groupDays {
		if groupDay == dayName {
			isGroupDay = true
			break
		}
	}

	return isGroupDay
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

func CalculateMoneyForLesson(db *sql.DB, price *float64, studentId string, groupId string, attendDate string, discountAmount, courseP, fixedSum *float64) error {
	var coursePrice float64
	err := db.QueryRow(`SELECT price FROM courses c join groups g on c.id=g.course_id where g.id=$1`, groupId).Scan(&coursePrice)
	if err != nil {
		return err
	}
	if fixedSum != nil {
		if discountAmount != nil {
			percent := *discountAmount * 100 / coursePrice
			fmt.Printf("bu yerda foizi %v", percent)
			teacherAmount := *fixedSum * percent / 100
			fmt.Printf("bu yerda teacher uchun amount %v", teacherAmount)
			*fixedSum = coursePrice - *discountAmount
			fmt.Printf("bu yerda fixed summa %v", *fixedSum)
			coursePrice = teacherAmount
		} else {
			coursePrice = *fixedSum
		}
	} else {
		if discountAmount != nil {
			coursePrice = coursePrice - *discountAmount
		}
	}

	*courseP = coursePrice
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

func NewTimoutContext(ctx context.Context, companyId string) (context.Context, context.CancelFunc) {
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, cancelFunc := context.WithTimeout(ctx, time.Second*60)
	return res, cancelFunc
}
