package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetAllStudent(condition string, page string, size string) (*pb.GetAllStudentResponse, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("invalid page value: %v", err)
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return nil, fmt.Errorf("invalid size value: %v", err)
	}

	offset := (pageInt - 1) * sizeInt

	countQuery := `SELECT COUNT(*) FROM students WHERE condition = $1`
	var totalCount int32
	err = r.db.QueryRow(countQuery, condition).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	query := `
    SELECT 
        s.id, s.name, s.gender, s.date_of_birth, s.phone, s.address, s.passport_id, s.additional_contact, 
        s.balance, s.condition, s.telegram_username, s.created_at,
        g.id AS group_id, g.name AS group_name, g.start_date, g.end_date, g.days, g.start_time,
        c.id AS course_id, c.title AS course_title, c.duration_lesson, c.course_duration, c.price,
        'exampleteachername' AS teacher_name,
        gs.condition AS student_group_condition, g.room_id, gs.last_specific_date AS student_activated_at
    FROM students s
    LEFT JOIN group_students gs ON s.id = gs.student_id
    LEFT JOIN groups g ON gs.group_id = g.id
    LEFT JOIN courses c ON g.course_id = c.id
    LEFT JOIN rooms r ON g.room_id = r.id
    WHERE s.condition = $1
    LIMIT $2 OFFSET $3;
    `

	rows, err := r.db.Query(query, condition, sizeInt, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var response pb.GetAllStudentResponse

	for rows.Next() {
		var student pb.GetGroupsAbsForStudent
		var group pb.GroupGetAllStudentAbs
		var course pb.AbsCourse

		var nullableGroupID, nullableGroupName, nullableGroupStartDate, nullableGroupEndDate sql.NullString
		var nullableDays pq.StringArray
		var nullableLessonStartTime sql.NullString
		var nullableCourseID, nullableCourseName sql.NullString
		var nullableLessonDuration, nullableCourseDuration sql.NullInt32
		var nullablePrice sql.NullFloat64
		var nullableStudentCondition sql.NullString
		var nullableRoomID sql.NullInt32
		var nullableStudentActivatedAt sql.NullTime

		err := rows.Scan(
			&student.Id, &student.Name, &student.Gender, &student.DateOfBirth, &student.Phone,
			&student.Address, &student.PassportId, &student.AdditionalContact, &student.Balance,
			&student.Condition, &student.TelegramUsername, &student.CreatedAt,
			&nullableGroupID, &nullableGroupName, &nullableGroupStartDate, &nullableGroupEndDate, &nullableDays, &nullableLessonStartTime,
			&nullableCourseID, &nullableCourseName, &nullableLessonDuration, &nullableCourseDuration, &nullablePrice,
			&group.TeacherName, &nullableStudentCondition, &nullableRoomID, &nullableStudentActivatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		if nullableGroupID.Valid {
			group.Id = nullableGroupID.String
			group.Name = nullableGroupName.String
			group.GroupStartDate = nullableGroupStartDate.String
			group.GroupEndDate = nullableGroupEndDate.String
			group.Days = nullableDays
			group.LessonStartTime = nullableLessonStartTime.String
			group.StudentCondition = nullableStudentCondition.String
			if nullableRoomID.Valid {
				group.RoomId = int32(nullableRoomID.Int32)
			}
			if nullableStudentActivatedAt.Valid {
				group.StudentActivatedAt = nullableStudentActivatedAt.Time.Format(time.RFC3339)
			}

			if nullableCourseID.Valid {
				course.Id = nullableCourseID.String
				course.Name = nullableCourseName.String
				course.LessonDuration = nullableLessonDuration.Int32
				course.CourseDuration = nullableCourseDuration.Int32
				course.Price = nullablePrice.Float64
				group.Course = &course
			}

			student.Groups = append(student.Groups, &group)
		}

		response.Response = append(response.Response, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	response.TotalCount = totalCount

	totalPages := (totalCount + int32(sizeInt) - 1) / int32(sizeInt)
	remainingPages := totalPages - int32(pageInt)
	if remainingPages < 0 {
		remainingPages = 0
	}
	response.TotalCount = remainingPages

	return &response, nil
}
func (r *StudentRepository) CreateStudent(createdBy string, phoneNumber string, name string, groupId string, address string, additionalContact string, dateFrom string, birthDate string, gender bool, passportId string, telegramUsername string) error {
	studentId := uuid.New()
	_, err := r.db.Exec(`INSERT INTO students(id, name, phone, date_of_birth, gender, telegram_username, passport_id, additional_contact, address) values ($1, $2,$3,$4,$5,$6,$7,$8,$9)`, studentId, name, phoneNumber, birthDate, gender, telegramUsername, passportId, additionalContact, address)
	if err != nil {
		return err
	}
	if groupId != "" && dateFrom != "" && createdBy != "" {
		_, err = r.db.Exec(`INSERT INTO group_students(id, group_id, student_id, created_by) values ($1 ,$2 ,$3 ,$4)`, uuid.New(), groupId, studentId, createdBy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *StudentRepository) UpdateStudent(studentId string, number string, name string, address string, additionalContact string, birth string, gender bool, passportId string) error {
	return nil
}

func (r *StudentRepository) DeleteStudent(studentId string) error {
	return nil
}

func (r *StudentRepository) AddToGroup(groupId string, studentIds []string, createdDate string) error {
	return nil
}
