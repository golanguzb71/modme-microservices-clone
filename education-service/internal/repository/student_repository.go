package repository

import (
	"database/sql"
	"education-service/internal/utils"
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
	studentQuery := `
    SELECT id, name, gender, date_of_birth, phone, address, passport_id, additional_contact, 
           balance, condition, telegram_username, created_at
    FROM students
    WHERE condition = $1
    ORDER BY id
    LIMIT $2 OFFSET $3`
	studentRows, err := r.db.Query(studentQuery, condition, sizeInt, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute student query: %v", err)
	}
	defer studentRows.Close()

	var students []*pb.GetGroupsAbsForStudent
	for studentRows.Next() {
		var student pb.GetGroupsAbsForStudent
		err := studentRows.Scan(
			&student.Id, &student.Name, &student.Gender, &student.DateOfBirth, &student.Phone,
			&student.Address, &student.PassportId, &student.AdditionalContact, &student.Balance,
			&student.Condition, &student.TelegramUsername, &student.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student row: %v", err)
		}
		students = append(students, &student)
	}
	groupQuery := `
    SELECT 
        s.id AS student_id,
        g.id AS group_id, g.name AS group_name, g.start_date, g.end_date, g.date_type, g.days, g.start_time,
        c.id AS course_id, c.title AS course_title, c.duration_lesson, c.course_duration, c.price,
        'exampleteachername' AS teacher_name,
        gs.condition AS student_group_condition, g.room_id, gs.last_specific_date AS student_activated_at
    FROM students s
    JOIN group_students gs ON s.id = gs.student_id
    JOIN groups g ON gs.group_id = g.id
    JOIN courses c ON g.course_id = c.id
    WHERE s.id = ANY($1)
    ORDER BY s.id, g.id`

	studentIDs := make([]string, len(students))
	for i, student := range students {
		studentIDs[i] = student.Id
	}
	groupRows, err := r.db.Query(groupQuery, pq.Array(studentIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to execute group query: %v", err)
	}
	defer groupRows.Close()
	studentMap := make(map[string]*pb.GetGroupsAbsForStudent)
	for _, student := range students {
		studentMap[student.Id] = student
	}
	for groupRows.Next() {
		var studentID string
		var group pb.GroupGetAllStudentAbs
		var course pb.AbsCourse

		err := groupRows.Scan(
			&studentID,
			&group.Id, &group.Name, &group.GroupStartDate, &group.GroupEndDate, &group.Type, pq.Array(&group.Days), &group.LessonStartTime,
			&course.Id, &course.Name, &course.LessonDuration, &course.CourseDuration, &course.Price,
			&group.TeacherName, &group.StudentCondition, &group.RoomId, &group.StudentActivatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group row: %v", err)
		}

		group.Course = &course
		student := studentMap[studentID]
		student.Groups = append(student.Groups, &group)
	}
	var response pb.GetAllStudentResponse
	for _, student := range students {
		response.Response = append(response.Response, student)
	}
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
	_, err := r.db.Exec(`UPDATE students SET phone =$1, name=$2, address =$3, additional_contact =$4, date_of_birth =$5, gender =$6, passport_id=$7 where id=$8`, number, name, address, additionalContact, birth, gender, passportId, studentId)
	if err != nil {
		return err
	}
	return nil
}
func (r *StudentRepository) DeleteStudent(studentId string) error {
	var cond string
	if err := r.db.QueryRow(`select condition from students where id = $1`, studentId).Scan(&cond); err != nil {
		return err
	}

	if cond == "ACTIVE" {
		_, err := r.db.Exec(`UPDATE students SET condition='ARCHIVED' where id=$1`, studentId)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.Exec(`UPDATE students SET condition='ACTIVE' where id=$1`, studentId)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *StudentRepository) AddToGroup(groupId string, studentIds []string, createdDate, createdBy string) error {
	query := `INSERT INTO group_students(id, group_id, student_id, condition, last_specific_date, created_by) values ($1 ,$2 ,$3 ,$4 , $5 , $6)`
	for _, data := range studentIds {
		_, err := r.db.Exec(query, uuid.New(), groupId, data, "FREEZE", createdDate, createdBy)
		if err != nil {
			continue
		}
	}
	return nil
}
func (r *StudentRepository) GetStudentById(id string) (*pb.GetStudentByIdResponse, error) {
	var result pb.GetStudentByIdResponse

	err := r.db.QueryRow(`SELECT id, name, gender, date_of_birth, phone, balance, created_at , condition
                          FROM students WHERE id = $1`, id).
		Scan(&result.Id, &result.Name, &result.Gender, &result.DateOfBirth, &result.Phone, &result.Balance, &result.CreatedAt, &result.Condition)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
        SELECT gs.created_at, gs.group_id, gs.condition, gs.last_specific_date, 
               g.name, g.date_type, g.days, g.start_time, g.start_date, g.end_date,
               r.id, r.capacity, r.title, 
               c.id, c.title, c.duration_lesson, c.course_duration, c.price, c.description , 'Shokruh' as teacher_name
        FROM group_students gs
        JOIN groups g ON g.id = gs.group_id
        JOIN rooms r ON g.room_id = r.id
        JOIN courses c ON c.id = g.course_id
        WHERE gs.student_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupStudent pb.GetGroupStudent
		var room pb.AbsRoom
		var course pb.AbsCourse

		err := rows.Scan(&groupStudent.StudentAddedAt, &groupStudent.Id, &groupStudent.StudentCondition, &groupStudent.StudentActivatedAt,
			&groupStudent.Name, &groupStudent.DateType, pq.Array(&groupStudent.Days), &groupStudent.LessonStartTime,
			&groupStudent.GroupStartDate, &groupStudent.GroupEndDate,
			&room.Id, &room.Capacity, &room.Name,
			&course.Id, &course.Name, &course.LessonDuration, &course.CourseDuration, &course.Price, &course.Description, &groupStudent.TeacherName)
		if err != nil {
			return nil, err
		}
		groupStudent.PriceForStudent = course.Price
		groupStudent.Room = &room
		groupStudent.Course = &course
		result.Groups = append(result.Groups, &groupStudent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}
func (r *StudentRepository) GetNoteByStudent(id string) (*pb.GetNotesByStudent, error) {
	rows, err := r.db.Query(`SELECT id, comment, created_at FROM student_note where student_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []*pb.AbsNote
	for rows.Next() {
		var note pb.AbsNote
		err = rows.Scan(&note.Id, &note.Comment, &note.CreatedAt)
		if err != nil {
			continue
		}
		notes = append(notes, &note)
	}
	return &pb.GetNotesByStudent{Notes: notes}, nil
}
func (r *StudentRepository) CreateNoteForStudent(note string, studentId string) (*pb.AbsResponse, error) {
	_, err := r.db.Exec(`INSERT INTO student_note(id , student_id, comment) values ($1,$2,$3)`, uuid.New(), studentId, note)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "created",
	}, nil
}
func (r *StudentRepository) DeleteStudentNote(id string) (*pb.AbsResponse, error) {
	_, err := r.db.Exec(`DELETE FROM student_note WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "deleted",
	}, nil
}
func (r *StudentRepository) SearchStudent(value string) (*pb.SearchStudentResponse, error) {
	query := `
        SELECT id, name, phone 
        FROM students 
        WHERE name ILIKE $1 OR phone ILIKE $2;
    `

	rows, err := r.db.Query(query, "%"+value+"%", "%"+value+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*pb.AbsStudent
	for rows.Next() {
		var student pb.AbsStudent
		if err := rows.Scan(&student.Id, &student.Name, &student.PhoneNumber); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.SearchStudentResponse{Students: students}, nil
}
func (r *StudentRepository) GetHistoryGroupById(id string) (*pb.GetHistoryGroupResponse, error) {
	var response pb.GetHistoryGroupResponse
	groupHistoryQuery := `
        SELECT 
            id,
            description,
            created_at
        FROM group_history
        WHERE group_id = $1
        ORDER BY created_at DESC`
	rows, err := r.db.Query(groupHistoryQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying group history: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var history pb.AbsHistory
		var createdAt time.Time
		err := rows.Scan(
			&history.Id,
			&history.Description,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning group history row: %v", err)
		}
		history.CreatedAt = createdAt.Format(time.RFC3339)
		response.History = append(response.History, &history)
	}
	studentHistoryQuery := `
        SELECT 
            s.id as student_id,
            s.name as student_name,
            g.id as group_id,
            g.name as group_name,
            g.teacher_id,
            g.start_date,
            g.end_date,
            g.start_time,
            g.date_type,
            gsch.condition,
            gsch.specific_date,
            gsch.created_at
        FROM group_student_condition_history gsch
        JOIN students s ON s.id = gsch.student_id
        JOIN groups g ON g.id = gsch.group_id
        WHERE gsch.group_id = $1
        ORDER BY gsch.created_at DESC`
	rows, err = r.db.Query(studentHistoryQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying student history: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var history pb.AbsStudentHistory
		var student pb.AbsStudent
		var group pb.AbsGroup
		var teacherId string
		var createdAt, specificDate time.Time

		err := rows.Scan(
			&student.Id,
			&student.Name,
			&group.Id,
			&group.Name,
			&teacherId,
			&group.GroupStartDate,
			&group.GroupEndDate,
			&group.LessonStartTime,
			&group.DateType,
			&history.Condition,
			&specificDate,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student history row: %v", err)
		}

		history.SpecificDate = specificDate.Format(time.RFC3339)
		history.CreatedAt = createdAt.Format(time.RFC3339)

		history.Student = &student
		history.Group = &group
		response.StudentHistory = append(response.StudentHistory, &history)
	}

	return &response, nil
}
func (r *StudentRepository) GetHistoryStudentById(id string) (*pb.GetHistoryStudentResponse, error) {
	var response pb.GetHistoryStudentResponse

	// Get student history
	studentHistoryQuery := `
        SELECT 
            id,
            description,
            created_at
        FROM student_history
        WHERE student_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.Query(studentHistoryQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying student history: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var history pb.AbsHistory
		var createdAt time.Time

		err := rows.Scan(
			&history.Id,
			&history.Description,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student history row: %v", err)
		}
		history.CreatedAt = createdAt.Format(time.RFC3339)
		response.History = append(response.History, &history)
	}

	// Get group condition history for the student
	groupHistoryQuery := `
        SELECT 
            s.id as student_id,
            s.name as student_name,
            g.id as group_id,
            g.name as group_name,
            g.teacher_id,
            g.start_date,
            g.end_date,
            g.start_time,
            g.date_type,
            gsch.condition,
            gsch.specific_date,
            gsch.created_at
        FROM group_student_condition_history gsch
        JOIN students s ON s.id = gsch.student_id
        JOIN groups g ON g.id = gsch.group_id
        WHERE gsch.student_id = $1
        ORDER BY gsch.created_at DESC`

	rows, err = r.db.Query(groupHistoryQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying group history: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var history pb.AbsStudentHistory
		var student pb.AbsStudent
		var group pb.AbsGroup
		var teacherId string
		var createdAt, specificDate time.Time

		err := rows.Scan(
			&student.Id,
			&student.Name,
			&group.Id,
			&group.Name,
			&teacherId,
			&group.GroupStartDate,
			&group.GroupEndDate,
			&group.LessonStartTime,
			&group.DateType,
			&history.Condition,
			&specificDate,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning group history row: %v", err)
		}

		// Format dates
		group.GroupStartDate = group.GroupStartDate
		group.GroupEndDate = group.GroupEndDate
		history.SpecificDate = specificDate.Format(time.RFC3339)
		history.CreatedAt = createdAt.Format(time.RFC3339)

		history.Student = &student
		history.Group = &group
		response.StudentHistory = append(response.StudentHistory, &history)
	}

	return &response, nil
}

func (r *StudentRepository) TransferLessonDate(groupId string, from string, to string) (*pb.AbsResponse, error) {
	validDay, err := utils.IsValidLessonDay(r.db, groupId, from)
	if err != nil {
		return nil, err
	}

	if !validDay {
		return &pb.AbsResponse{
			Status:  403,
			Message: "The selected 'from' date does not match the group's lesson days",
		}, nil
	}
	var checker bool
	err = r.db.QueryRow(`SELECT exists(SELECT 1 FROM transfer_lesson where group_id=$1 and real_date=$2 and transfer_date=$3)`, groupId, from, to).Scan(&checker)
	if err != nil {
		return nil, err
	}
	if checker {
		_, err = r.db.Exec(`DELETE FROM transfer_lesson where group_id=$1 and real_date=$2 and transfer_date=$3`, groupId, from, to)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = r.db.Exec(`INSERT INTO transfer_lesson(id, group_id, real_date, transfer_date) values ($1, $2, $3, $4)`, uuid.New(), groupId, from, to)
		if err != nil {
			return nil, err
		}
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "accomplished",
	}, nil
}
