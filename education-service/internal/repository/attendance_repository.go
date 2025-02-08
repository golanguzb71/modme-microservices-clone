package repository

import (
	"context"
	"database/sql"
	"education-service/internal/clients"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
	"time"
)

type AttendanceRepository struct {
	db                *sql.DB
	financeClient     *clients.FinanceClient
	financeClientChan chan *clients.FinanceClient
}

type Groups struct {
	Id                  string
	Name                string
	CourseId            int
	DateType            string
	Days                []string
	StartTime           string
	StartDate           string
	EndDate             string
	IsArchived          bool
	LessonCountOnPeriod int32
	CreatedAt           string
}
type Attendance struct {
	IsDiscounted  bool
	DiscountOwner string
	Price         float32
	GroupId       int64
	StudentId     string
	StudentName   string
	TeacherId     string
	AttendDate    string
	Status        int
	CreatedAt     time.Time
	CreatedBy     string
	CreatorRole   string
	PriceType     string
	TotalCount    float64
	CoursePrice   float64
}

func (r *AttendanceRepository) ensureFinanceClient() error {
	if r.financeClient == nil {
		select {
		case client := <-r.financeClientChan:
			r.financeClient = client
		case <-time.After(5 * time.Second):
			return fmt.Errorf("failed to initialize GroupClient within timeout")
		}
	}
	return nil
}

func NewAttendanceRepository(db *sql.DB, financeClientChan chan *clients.FinanceClient) *AttendanceRepository {
	return &AttendanceRepository{db: db, financeClientChan: financeClientChan}
}
func (r *AttendanceRepository) CreateAttendance(ctx context.Context, companyId, groupId string, studentId string, teacherId string, attendDate string, status int32, actionById, actionByRole string) error {
	if err := r.ensureFinanceClient(); err != nil {
		return fmt.Errorf("error while ensuring finance client %v", err)
	}

	var (
		isDiscounted bool
		price        float64
		priceType    string
		totalCount   int
		coursePrice  float64
	)
	ctx, c := utils.NewTimoutContext(ctx, companyId)
	defer c()
	resp, err := r.financeClient.GetTeacherSalaryByTeacherID(ctx, teacherId)
	if err != nil {
		return errors.New("error while getting teacher salary information")
	}
	if !utils.CheckGroupAndTeacher(r.db, groupId, "TEACHER", teacherId) {
		return fmt.Errorf("oops this teacherid not the same for this group")
	}
	discountAmount, discountOwner := r.financeClient.GetDiscountByStudentId(ctx, studentId, groupId)

	if resp.Type == "FIXED" {
		priceType = "FIXED"
		totalCount = int(resp.Amount)
		f := float64(totalCount)
		if discountAmount != nil {
			isDiscounted = true
			priceType = "FIXED_DISCOUNT"
		}
		if err = utils.CalculateMoneyForLesson(r.db, &price, studentId, groupId, attendDate, discountAmount, &coursePrice, &f); err != nil {
			return errors.New("error while getting calculate money")
		}
	} else {
		priceType = "PERCENT"
		if discountAmount != nil {
			isDiscounted = true
			priceType = "PERCENT_DISCOUNT"
			totalCount = int(resp.Amount)
		}
		if err = utils.CalculateMoneyForLesson(r.db, &price, studentId, groupId, attendDate, discountAmount, &coursePrice, nil); err != nil {
			return errors.New("error while getting calculate money")
		}
	}
	query := `
     	INSERT INTO attendance (is_discounted, discount_owner,  price , group_id , student_id , teacher_id, attend_date, status , created_at , created_by , creator_role , company_id , price_type , total_count , course_price)
        VALUES ($1, $2, $3, $4, $5 , $6 , $7 , $8, $9 , $10 , $11 , $12 , $13 , $14, $15)
        ON CONFLICT DO NOTHING
    `
	_, err = r.db.Exec(query, isDiscounted, discountOwner, price, groupId, studentId, teacherId, attendDate, status, time.Now(), actionById, actionByRole, companyId, priceType, totalCount, coursePrice)
	if err != nil {
		return fmt.Errorf("error while creating attendance %v", err)
	}
	return err
}
func (r *AttendanceRepository) DeleteAttendance(groupId string, studentId string, teacherId string, attendDate string) error {
	if !utils.CheckGroupAndTeacher(r.db, groupId, "TEACHER", teacherId) {
		return fmt.Errorf("oops this teacherid not the same for this group")
	}
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
	return nil
}
func (r *AttendanceRepository) GetAttendanceByGroupAndDateRange(companyId string, ctx context.Context, groupId string, fromDate time.Time, tillDate time.Time, withOutdated bool, actionRole, actionId string) (*pb.GetAttendanceResponse, error) {
	if !utils.CheckGroupAndTeacher(r.db, groupId, actionRole, actionId) {
		return nil, status.Errorf(codes.Aborted, "Ooops. this group not found in your groupList")
	}

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
func (r *AttendanceRepository) IsHaveTransferredLesson(groupID string) bool {
	query := `
		SELECT COUNT(*) 
		FROM transfer_lesson 
		WHERE group_id = $1 AND transfer_date = CURRENT_DATE`

	var count int
	err := r.db.QueryRow(query, groupID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}
func (r *AttendanceRepository) GetAllGroupsByTeacherId(teacherId, from, to string) []Groups {
	var response []Groups
	rows, err := r.db.Query(`SELECT id,
       name,
       course_id,
       date_type,
       days,
       start_time,
       start_date,
       end_date,
       is_archived,
       created_at FROM groups where teacher_id=$1`, teacherId)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var group Groups
		err := rows.Scan(&group.Id, &group.Name, &group.CourseId, &group.DateType, pq.Array(&group.Days), &group.StartTime, &group.StartDate, &group.EndDate, &group.IsArchived, &group.CreatedAt)
		if err != nil {
			return nil
		}
		group.LessonCountOnPeriod = r.lessonCounter(from, to, group.Id)
		response = append(response, group)
	}

	return response
}

func (r *AttendanceRepository) lessonCounter(from, to, groupId string) int32 {
	query := `
        WITH RECURSIVE dates AS (
            SELECT $1::date AS date
            UNION ALL
            SELECT date + 1
            FROM dates
            WHERE date < $2::date
        ),
        group_dates AS (
            SELECT DISTINCT d.date
            FROM dates d
            JOIN groups g ON g.id = $3::bigint
            WHERE
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
        SELECT COUNT(*)
        FROM group_dates;
    `

	var lessonCount int32
	err := r.db.QueryRow(query, from, to, groupId).Scan(&lessonCount)
	if err != nil {
		return 0
	}

	return lessonCount
}

func (r *AttendanceRepository) GetAttendanceByTeacherAndGroup(companyId, teacherId string, groupId string, from string, to string) (map[string][]Attendance, error) {
	query := `
		SELECT 
			is_discounted,
			discount_owner,
			price,
			group_id,
			student_id,
			s.name,
			teacher_id,
			attend_date,
			status,
			a.created_at,
			created_by,
			creator_role,
			price_type,
			total_count,
			course_price
		FROM 
			attendance a
		JOIN students s
		ON a.student_id=s.id
		WHERE 
			teacher_id = $1 AND 
			group_id = $2 AND 
			attend_date BETWEEN $3 AND $4
			AND
		    creator_role='TEACHER'
			AND a.company_id=$5
	`

	rows, err := r.db.Query(query, teacherId, groupId, from, to, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendanceMap := make(map[string][]Attendance)

	for rows.Next() {
		var attendance Attendance
		err := rows.Scan(
			&attendance.IsDiscounted,
			&attendance.DiscountOwner,
			&attendance.Price,
			&attendance.GroupId,
			&attendance.StudentId,
			&attendance.StudentName,
			&attendance.TeacherId,
			&attendance.AttendDate,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.CreatedBy,
			&attendance.CreatorRole,
			&attendance.PriceType,
			&attendance.TotalCount,
			&attendance.CoursePrice,
		)
		if err != nil {
			return nil, err
		}

		attendanceMap[attendance.StudentId] = append(attendanceMap[attendance.StudentId], attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attendanceMap, nil
}
