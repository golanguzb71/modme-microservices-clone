package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strconv"
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
        g.id AS group_id, g.name AS group_name, g.start_date, g.end_date, g.days, g.start_time,
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
			&group.Id, &group.Name, &group.GroupStartDate, &group.GroupEndDate, pq.Array(&group.Days), &group.LessonStartTime,
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
	err := r.db.QueryRow(`SELECT id, name, gender, date_of_birth, phone, balance, created_at FROM students where id=$1`, id).Scan(&result.Id, &result.Name, &result.Gender, &result.DateOfBirth, &result.Phone, &result.Balance, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`SELECT group_id, condition, last_specific_date FROM group_students where student_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groupStudent pb.GetGroupStudent
	for rows.Next() {
		//		`SELECT group_id,
		//       condition,
		//       last_specific_date,
		//       g.name,
		//       g.date_type,
		//       g.days,
		//       r.id,
		//       r.capacity,
		//       r.title,
		//       c.id,
		//       c.title,
		//       c.duration_lesson,
		//       c.course_duration,
		//       c.price,
		//       c.description,
		//       g.start_time,
		//       g.start_date,
		//       g.end_date,
		//       condition,
		//       last_specific_date,
		//       c.price
		//FROM group_students gs
		//         join groups g on g.id = gs.group_id
		//         join rooms r on g.room_id = r.id
		//         join courses c on c.id = g.course_id
		//where student_id = '20b1926b-a3d9-417f-8662-97ee6aa43618'`
		rows.Scan(&groupStudent.Id, &groupStudent.StudentCondition, &groupStudent.StudentActivatedAt)
	}
}
func (r *StudentRepository) GetNoteByStudent(id string) (*pb.GetNotesByStudent, error) {

}
func (r *StudentRepository) CreateNoteForStudent(note string, studentId string) (*pb.AbsResponse, error) {

}
func (r *StudentRepository) DeleteStudentNote(id string) (*pb.AbsResponse, error) {

}
func (r *StudentRepository) SearchStudent(value string) (*pb.SearchStudentResponse, error) {

}
func (r *StudentRepository) GetHistoryGroupById(id string) (*pb.GetHistoryGroupResponse, error) {

}
func (r *StudentRepository) GetHistoryStudentById(id string) (*pb.GetHistoryStudentResponse, error) {

}
func (r *StudentRepository) TransferLessonDate(groupId string, from string, to string) (*pb.AbsResponse, error) {

}
