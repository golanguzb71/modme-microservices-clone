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
	var groupEndDate string
	var groupDays []string
	var dateType string

	query := `
        SELECT c.price, c.course_duration, g.start_date, g.end_date, g.days, g.date_type
        FROM groups g
        JOIN courses c ON g.course_id = c.id
        WHERE g.id = $1
    `

	err := db.QueryRow(query, groupId).Scan(&coursePrice, &courseDurationLesson, &groupStartDate, &groupEndDate, pq.Array(&groupDays), &dateType)
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

	startOfMonth := time.Date(tillDateParsed.Year(), tillDateParsed.Month(), 1, 0, 0, 0, 0, tillDateParsed.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Parse the timestamp format from the database
	groupEndDateParsed, err := time.Parse(time.RFC3339, groupEndDate)
	if err != nil {
		// Try alternate format if RFC3339 fails
		groupEndDateParsed, err = time.Parse("2006-01-02T15:04:05Z", groupEndDate)
		if err != nil {
			return 0, fmt.Errorf("error parsing group end date: %v", err)
		}
	}

	// Calculate total days in month
	daysInMonth := endOfMonth.Day()

	// Calculate effective end date and days for price calculation
	effectiveEndDate := endOfMonth
	if groupEndDateParsed.Before(endOfMonth) {
		effectiveEndDate = groupEndDateParsed
	}

	// Calculate the proportional price based on days
	daysToConsider := effectiveEndDate.Day()
	proportionalPrice := (coursePrice * float64(daysToConsider)) / float64(daysInMonth)

	// Get lesson dates for the actual period
	lessonDates := getLessonDatesInMonth(groupDays, dateType, startOfMonth, effectiveEndDate)
	if len(lessonDates) == 0 {
		return 0, fmt.Errorf("no lessons scheduled for the current month")
	}

	firstLessonDate := lessonDates[0]

	// If we're before or on the first lesson date, return proportional price
	if !tillDateParsed.After(firstLessonDate) {
		return math.Round(proportionalPrice), nil
	}

	// Count passed lessons
	passedLessons := 0
	for _, lessonDate := range lessonDates {
		if lessonDate.Before(tillDateParsed) {
			passedLessons++
		}
	}

	// Calculate remaining money based on proportional price
	pricePerLesson := proportionalPrice / float64(len(lessonDates))
	remainingMoney := proportionalPrice - (float64(passedLessons) * pricePerLesson)

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
	err := db.QueryRow(`SELECT price FROM courses c JOIN groups g ON c.id = g.course_id WHERE g.id = $1`, groupId).Scan(&coursePrice)
	if err != nil {
		return fmt.Errorf("failed to fetch course price: %v", err)
	}

	if fixedSum != nil {
		if discountAmount != nil {
			fmt.Println("discount amount ", *discountAmount)
			percent := (*discountAmount * 100) / coursePrice
			teacherAmount := (*fixedSum * percent) / 100
			*fixedSum = coursePrice - *discountAmount
			fmt.Println("fixed summa narxi ", *fixedSum)
			coursePrice = teacherAmount
		} else {
			fmt.Println("discount yoq ekan ")
			coursePrice = *fixedSum
		}
	} else if discountAmount != nil {
		fmt.Println("fixed summada emas ekan")
		coursePrice -= *discountAmount
	}

	fmt.Println("calculate moneyga kridi")
	*courseP = coursePrice

	parsedDate, err := time.Parse("2006-01-02", attendDate)
	if err != nil {
		return fmt.Errorf("invalid attendDate format: %v", err)
	}

	firstOfMonth := time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, parsedDate.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	var lessonCount int
	query := `
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
    `
	err = db.QueryRow(query, firstOfMonth, lastOfMonth, groupId).Scan(&lessonCount)
	if err != nil {
		return fmt.Errorf("failed to count lesson days: %v", err)
	}

	// Handle case where no lessons are found
	if lessonCount == 0 {
		return fmt.Errorf("no lessons found in the month for group %s", groupId)
	}

	// Calculate and set the price per lesson
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
