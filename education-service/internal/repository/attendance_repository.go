package repository

import (
	"context"
	"database/sql"
	"education-service/proto/pb"
	"fmt"
	"sort"
	"time"
)

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) CreateAttendance(groupId string, studentId string, teacherId string, attendDate string, status int32) error {
	query := `
        INSERT INTO attendance (group_id, student_id, teacher_id, attend_date, status)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT DO NOTHING 
    `
	_, err := r.db.Exec(query, groupId, studentId, teacherId, attendDate, status)
	return err
}
func (r *AttendanceRepository) DeleteAttendance(groupId string, studentId string, teacherId string, attendDate string) error {
	query := `
        DELETE FROM attendance
        WHERE group_id = $1
          AND student_id = $2
          AND teacher_id = $3
          AND attend_date = $4
    `
	result, err := r.db.Exec(query, groupId, studentId, teacherId, attendDate)
	if err != nil {
		return fmt.Errorf("failed to delete attendance: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("attendance record not found for group_id: %s, student_id: %s, teacher_id: %s, attend_date: %s", groupId, studentId, teacherId, attendDate)
	}
	// writing teacher finance remove
	return nil
}
func (r *AttendanceRepository) GetAttendanceByGroupAndDateRange(ctx context.Context, groupId string, fromDate time.Time, tillDate time.Time, withOutdated bool) (*pb.GetAttendanceResponse, error) {
	response := &pb.GetAttendanceResponse{
		Days:     make([]*pb.Day, 0),
		Students: make([]*pb.Student, 0),
	}

	daysQuery := `
        WITH RECURSIVE dates AS (
            SELECT $1::date AS date
            UNION ALL
            SELECT date + 1
            FROM dates
            WHERE date < $2::date
        ),
        group_dates AS (
            SELECT DISTINCT d.date::text, tl.transfer_date::text
            FROM dates d
            JOIN groups g ON g.id = $3::bigint
            LEFT JOIN transfer_lesson tl ON tl.group_id = g.id AND tl.real_date = d.date
            WHERE (
                (d.date >= g.start_date AND d.date <= LEAST(g.end_date, $2::date))
                AND EXTRACT(DOW FROM d.date) = ANY(
                    SELECT CASE day
                        WHEN 'DUSHANBA' THEN 1
                        WHEN 'SESHANBA' THEN 2
                        WHEN 'CHORSHANBA' THEN 3
                        WHEN 'PAYSHANBA' THEN 4
                        WHEN 'JUMA' THEN 5
                        WHEN 'SHANBA' THEN 6
                        WHEN 'YAKSHANBA' THEN 0
                    END
                    FROM unnest(g.days) AS day
                )
            )
            ORDER BY d.date::text
        )
        SELECT date, transfer_date 
        FROM group_dates
        ORDER BY date, COALESCE(transfer_date, date)
    `

	rows, err := r.db.QueryContext(ctx, daysQuery, fromDate, tillDate, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateStr, transferDateStr sql.NullString
		if err := rows.Scan(&dateStr, &transferDateStr); err != nil {
			return nil, err
		}
		response.Days = append(response.Days, &pb.Day{
			Date:         dateStr.String,
			TransferDate: transferDateStr.String,
		})
	}

	studentsQuery := `
        WITH student_attendance AS (
            SELECT 
                a.student_id,
                a.attend_date::text,
                a.status,
                a.teacher_id,
                a.created_at
            FROM attendance a
            WHERE a.group_id = $1::bigint
            AND a.attend_date BETWEEN $2::date AND $3::date
            ORDER BY a.attend_date, a.created_at
        ),
        last_activation AS (
            SELECT 
                student_id,
                created_at as activated_at
            FROM group_student_condition_history
            WHERE group_id = $1::bigint
            AND current_condition = 'ACTIVE'
            AND created_at = (
                SELECT MAX(created_at)
                FROM group_student_condition_history gsch2
                WHERE gsch2.student_id = group_student_condition_history.student_id
                AND gsch2.group_id = group_student_condition_history.group_id
                AND gsch2.current_condition = 'ACTIVE'
            )
        )
        SELECT 
            gs.student_id,
            gs.created_at as added_at,
            gs.condition,
            sa.attend_date,
            sa.status,
            sa.teacher_id,
            sa.created_at as attendance_created_at,
            s.name,
            s.phone,
            s.date_of_birth,
            s.gender,
            s.balance,
            s.created_at as student_created_at,
            s.condition as student_condition,
            la.activated_at
        FROM group_students gs
        LEFT JOIN student_attendance sa ON gs.student_id = sa.student_id
        LEFT JOIN students s ON gs.student_id = s.id
        LEFT JOIN last_activation la ON gs.student_id = la.student_id
        WHERE gs.group_id = $1::bigint
        AND ($4 OR gs.condition = 'ACTIVE')
        ORDER BY gs.created_at, sa.attend_date, sa.created_at
    `

	rows, err = r.db.QueryContext(ctx, studentsQuery, groupId, fromDate, tillDate, withOutdated)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	studentMap := make(map[string]*pb.Student)

	for rows.Next() {
		var (
			studentId, name, phone, condition, studentCondition      string
			dateOfBirth, addedAt, createdAt, activatedAt, attendDate sql.NullString
			teacherId                                                sql.NullString
			status                                                   sql.NullInt32
			gender                                                   sql.NullBool
			balance                                                  float64
			attendanceCreatedAt                                      sql.NullTime
		)

		if err := rows.Scan(
			&studentId,
			&addedAt,
			&condition,
			&attendDate,
			&status,
			&teacherId,
			&attendanceCreatedAt,
			&name,
			&phone,
			&dateOfBirth,
			&gender,
			&balance,
			&createdAt,
			&studentCondition,
			&activatedAt,
		); err != nil {
			return nil, err
		}

		student, exists := studentMap[studentId]
		if !exists {
			student = &pb.Student{
				Id:          studentId,
				Name:        name,
				Phone:       phone,
				DateOfBirth: dateOfBirth.String,
				Gender:      gender.Bool,
				Balance:     balance,
				CreatedAt:   createdAt.String,
				ActivatedAt: activatedAt.String,
				AddedAt:     addedAt.String,
				Condition:   condition,
				Attendance:  make([]*pb.Attendance, 0),
			}
			studentMap[studentId] = student
		}

		if attendDate.Valid && teacherId.Valid {
			attendance := &pb.Attendance{
				AttendDate: attendDate.String,
				IsCome:     status.Int32 == 1,
				StudentId:  studentId,
				TeacherId:  teacherId.String,
			}
			student.Attendance = append(student.Attendance, attendance)
		}
	}

	students := make([]*pb.Student, 0, len(studentMap))
	for _, student := range studentMap {
		students = append(students, student)
	}

	sort.Slice(students, func(i, j int) bool {
		return students[i].AddedAt < students[j].AddedAt
	})

	for _, student := range students {
		sort.Slice(student.Attendance, func(i, j int) bool {
			return student.Attendance[i].AttendDate < student.Attendance[j].AttendDate
		})
	}
	response.Students = students
	return response, nil
}
func (r *AttendanceRepository) IsValidGroupDay(ctx context.Context, groupId string, date time.Time) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM groups g
            WHERE g.id = $1 AND
            $2::date BETWEEN g.start_date AND g.end_date AND
            EXTRACT(DOW FROM $2::date) = ANY(
                SELECT CASE day
                    WHEN 'DUSHANBA' THEN 1
                    WHEN 'SESHANBA' THEN 2
                    WHEN 'CHORSHANBA' THEN 3
                    WHEN 'PAYSHANBA' THEN 4
                    WHEN 'JUMA' THEN 5
                    WHEN 'SHANBA' THEN 6
                    WHEN 'YAKSHANBA' THEN 0
                END
                FROM unnest(g.days) AS day
            )
        )
    `
	var exists bool
	err := r.db.QueryRowContext(ctx, query, groupId, date).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
