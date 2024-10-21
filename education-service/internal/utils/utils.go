package utils

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
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
	for _, lessonDay := range lessonDays {
		if strings.ToUpper(lessonDay) == strings.ToUpper(dayOfWeek) {
			return true, nil
		}
	}

	return false, nil
}
