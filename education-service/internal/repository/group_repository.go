package repository

import (
	"database/sql"
	"education-service/internal/clients"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type GroupRepository struct {
	db         *sql.DB
	userClient *clients.UserClient
}

func NewGroupRepository(db *sql.DB, userClient *clients.UserClient) *GroupRepository {
	return &GroupRepository{db: db, userClient: userClient}
}
func (r *GroupRepository) CreateGroup(name string, courseId int32, teacherId string, dateType string, days []string, roomId int32, lessonStartTime string, groupStartDate string, groupEndDate string) (string, error) {
	query := `
		INSERT INTO groups(course_id, teacher_id, room_id, date_type, days, start_time, start_date, end_date, is_archived, name) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id`

	var groupId string
	err := r.db.QueryRow(query, courseId, teacherId, roomId, dateType, pq.Array(days), lessonStartTime, groupStartDate, groupEndDate, false, name).Scan(&groupId)
	if err != nil {
		return "", err
	}
	return groupId, nil
}
func (r *GroupRepository) UpdateGroup(id string, name string, courseId int32, teacherId string, dateType string, days []string, roomId int32, lessonStartTime string, groupStartDate string, groupEndDate string) error {
	query := `UPDATE groups SET course_id=$1, teacher_id=$2, room_id=$3, date_type=$4, days=$5, start_time=$6, start_date=$7, end_date=$8, name=$9 WHERE id=$10`
	_, err := r.db.Exec(query, courseId, teacherId, roomId, dateType, pq.Array(days), lessonStartTime, groupStartDate, groupEndDate, name, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupRepository) DeleteGroup(id string) error {
	query := `UPDATE groups SET is_archived = NOT is_archived WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupRepository) GetGroup(page, size int32, isArchive bool) (*pb.GetGroupsResponse, error) {
	offset := (page - 1) * size
	query := `SELECT g.id, g.course_id, COALESCE(c.title, 'Unknown Course') as course_title, 
       g.teacher_id,
       g.room_id, COALESCE(r.title, 'Unknown Room') as room_title, r.capacity, g.start_date, g.end_date, g.is_archived, 
       g.name, 
       COUNT(gs.id) as student_count, 
       g.created_at , g.days , g.start_time , g.date_type
FROM groups g
LEFT JOIN courses c ON g.course_id = c.id
LEFT JOIN rooms r ON g.room_id = r.id
LEFT JOIN group_students gs ON g.id = gs.group_id
WHERE g.is_archived = $1
GROUP BY g.id, c.title, r.title, r.capacity
LIMIT $2 OFFSET $3;`

	rows, err := r.db.Query(query, isArchive, size, offset)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var groups []*pb.GetGroupAbsResponse
	for rows.Next() {
		var group pb.GetGroupAbsResponse
		var studentCount int32
		var course pb.AbsCourse
		var room pb.AbsRoom

		err = rows.Scan(
			&group.Id, &course.Id, &course.Name,
			&group.TeacherId, &room.Id, &room.Name, &room.Capacity, &group.StartDate, &group.EndDate,
			&group.IsArchived, &group.Name, &studentCount, &group.CreatedAt, pq.Array(&group.Days), &group.LessonStartTime, &group.DateType,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		teacherName, err := r.userClient.GetTeacherById(group.TeacherId)
		if err != nil {
			continue
		}
		group.TeacherName = teacherName
		group.Course = &course
		group.Room = &room
		group.StudentCount = studentCount

		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	var totalCount int32
	countQuery := `SELECT COUNT(*) FROM groups WHERE is_archived = $1;`
	err = r.db.QueryRow(countQuery, isArchive).Scan(&totalCount)
	if err != nil {
		log.Printf("Error counting total groups: %v", err)
		return nil, fmt.Errorf("error counting total groups: %w", err)
	}

	totalPageCount := (totalCount + size - 1) / size

	return &pb.GetGroupsResponse{Groups: groups, TotalPageCount: totalPageCount}, nil
}
func (r *GroupRepository) GetGroupById(id, actionRole, actionId string) (*pb.GetGroupAbsResponse, error) {
	if !utils.CheckGroupAndTeacher(r.db, id, actionRole, actionId) {
		return nil, status.Errorf(codes.Aborted, "Ooops. this group not found in your groupList")
	}
	query := `SELECT g.id, g.course_id, c.title as course_title, 
              g.room_id, COALESCE(r.title, 'Unknown Room') as room_title, r.capacity, g.start_date, g.end_date, g.is_archived, g.name,
              COUNT(gs.id)  as student_count, 
              g.created_at , g.days , g.start_time , g.date_type , c.course_duration ,c.duration_lesson , c.description , c.price , g.teacher_id
              FROM groups g
              LEFT JOIN courses c ON g.course_id = c.id
              LEFT JOIN rooms r ON g.room_id = r.id
              LEFT JOIN group_students gs ON g.id = gs.group_id
              WHERE g.id = $1
              GROUP BY g.id, c.title, r.title, r.capacity , c.course_duration , c.duration_lesson , c.description , c.price`

	var group pb.GetGroupAbsResponse
	var studentCount sql.NullInt32
	var course pb.AbsCourse
	var room pb.AbsRoom

	err := r.db.QueryRow(query, id).Scan(
		&group.Id, &course.Id, &course.Name, &room.Id, &room.Name, &room.Capacity, &group.StartDate, &group.EndDate,
		&group.IsArchived, &group.Name, &studentCount, &group.CreatedAt, pq.Array(&group.Days), &group.LessonStartTime, &group.DateType, &course.CourseDuration, &course.LessonDuration, &course.Description, &course.Price, &group.TeacherId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group with id %s not found", id)
		}
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	teacherName, _ := r.userClient.GetTeacherById(group.TeacherId)
	group.TeacherName = teacherName
	group.Course = &course
	group.Room = &room
	group.StudentCount = studentCount.Int32
	return &group, nil
}
func (r *GroupRepository) GetGroupByCourseId(courseId string) (*pb.GetGroupsByCourseResponse, error) {
	query := `
        SELECT 
            g.teacher_id,
            g.id,
            g.start_date,
            g.end_date,
            g.date_type,
            g.start_time,
            g.name
        FROM groups g
        WHERE g.course_id = $1
    `

	rows, err := r.db.Query(query, courseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response pb.GetGroupsByCourseResponse
	for rows.Next() {
		var groupResponse pb.GetGroupByCourseAbsResponse
		var startDate, endDate, lessonStartTime, dateType, name, teacherId sql.NullString

		err = rows.Scan(
			&teacherId,
			&groupResponse.Id,
			&startDate,
			&endDate,
			&dateType,
			&lessonStartTime,
			&name,
		)
		if err != nil {
			return nil, err
		}
		teacherName, err := r.userClient.GetTeacherById(teacherId.String)
		if err != nil {
			return nil, err
		}
		groupResponse.TeacherName = teacherName
		groupResponse.GroupName = name.String
		groupResponse.GroupStartDate = startDate.String
		groupResponse.GroupEndDate = endDate.String
		groupResponse.DateType = dateType.String
		groupResponse.LessonStartTime = lessonStartTime.String
		response.Groups = append(response.Groups, &groupResponse)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(response.Groups) == 0 {
		return nil, fmt.Errorf("no groups found for course id: %s", courseId)
	}

	return &response, nil
}
func (r *GroupRepository) GetGroupByTeacherId(teacherId string, archived bool) (*pb.GetGroupsByTeacherResponse, error) {
	query := `
		SELECT 
		    g.id,
			g.name, 
			c.title AS course_name, 
			r.title AS room_name, 
			g.start_time, 
			g.date_type, 
			g.start_date, 
			g.end_date, 
			(
				SELECT COUNT(gs.student_id)
				FROM group_students gs
				WHERE gs.group_id = g.id
				AND gs.condition = 'ACTIVE'
			) AS active_student_count
		FROM groups g
		INNER JOIN courses c ON g.course_id = c.id
		LEFT JOIN rooms r ON g.room_id = r.id
		WHERE g.teacher_id = $1 AND g.is_archived = $2
	`
	rows, err := r.db.Query(query, teacherId, archived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response pb.GetGroupsByTeacherResponse
	for rows.Next() {
		var group pb.GetGroupByTeacherAbs
		var activeStudentCount int32
		if err := rows.Scan(&group.Id, &group.Name, &group.CourseName, &group.RoomName, &group.LessonStartTime, &group.DayType, &group.GroupStartAt, &group.GroupEndAt, &activeStudentCount); err != nil {
			return nil, err
		}
		group.ActiveStudentCount = activeStudentCount
		students, err := r.GetStudentsByGroupId(group.Id)
		if err != nil {
			return nil, err
		}
		group.Students = students
		response.Groups = append(response.Groups, &group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &response, nil
}
func (r *GroupRepository) GetStudentsByGroupId(groupId string) ([]*pb.AbsStudent, error) {
	query := `
		SELECT s.id, s.name, s.phone 
		FROM students s
		INNER JOIN group_students gs ON gs.student_id = s.id
		WHERE gs.group_id = $1 AND gs.condition = 'ACTIVE'
	`
	rows, err := r.db.Query(query, groupId)
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

	return students, nil
}
func (r *GroupRepository) GetCommonInformationEducation() (*pb.GetCommonInformationEducationResponse, error) {
	// Initialize the response to avoid nil pointer dereference
	response := new(pb.GetCommonInformationEducationResponse)

	var leaveGroupCount, activeGroupCount, activeStudentCount, debtorsCount int32

	// Query to get leaveGroupCount
	err := r.db.QueryRow(`SELECT COUNT(id) FROM group_students where condition='DELETE'`).Scan(&leaveGroupCount)
	if err != nil {
		leaveGroupCount = 0
	}

	// Query to get activeGroupCount
	err = r.db.QueryRow(`SELECT COUNT(id) FROM groups where is_archived=false`).Scan(&activeGroupCount)
	if err != nil {
		activeGroupCount = 0
	}

	// Query to get activeStudentCount
	err = r.db.QueryRow(`SELECT count(id) FROM students where condition='ACTIVE'`).Scan(&activeStudentCount)
	if err != nil {
		activeStudentCount = 0
	}

	// Query to get debtorsCount
	err = r.db.QueryRow(`SELECT COUNT(id) FROM students where balance < 0`).Scan(&debtorsCount)
	if err != nil {
		debtorsCount = 0
	}

	// Assign the values to the response fields
	response.DebtorsCount = debtorsCount
	response.LeaveGroupCount = leaveGroupCount
	response.ActiveGroupCount = activeGroupCount
	response.ActiveStudentCount = activeStudentCount

	return response, nil
}

func (r *GroupRepository) GetGroupsByStudentId(studentId string) (*pb.GetGroupsByStudentResponse, error) {
	var (
		comments []*pb.DebtorComment
		groups   []*pb.DebtorGroup
	)

	commentQuery := `
		SELECT id, comment
		FROM student_note 
		WHERE student_id = $1
	`
	commentRows, err := r.db.Query(commentQuery, studentId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %v", err)
	}
	defer commentRows.Close()

	for commentRows.Next() {
		var comment pb.DebtorComment
		if err := commentRows.Scan(&comment.CommentId, &comment.Comment); err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %v", err)
		}
		comments = append(comments, &comment)
	}

	groupQuery := `
		SELECT g.id, g.name AS course_title
		FROM group_students gs
		JOIN groups g ON gs.group_id = g.id
		WHERE gs.student_id = $1
	`
	groupRows, err := r.db.Query(groupQuery, studentId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups: %v", err)
	}
	defer groupRows.Close()

	for groupRows.Next() {
		var group pb.DebtorGroup
		if err := groupRows.Scan(&group.GroupId, &group.GroupName); err != nil {
			return nil, fmt.Errorf("failed to scan group row: %v", err)
		}
		groups = append(groups, &group)
	}

	return &pb.GetGroupsByStudentResponse{
		Comments: comments,
		Groups:   groups,
	}, nil
}
