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
	"log"
	"strconv"
)

type GroupRepository struct {
	db         *sql.DB
	userClient *clients.UserClient
}

func NewGroupRepository(db *sql.DB, userClient *clients.UserClient) *GroupRepository {
	return &GroupRepository{db: db, userClient: userClient}
}
func (r *GroupRepository) CreateGroup(companyId string, name string, courseId int32, teacherId string, dateType string, days []string, roomId int32, lessonStartTime string, groupStartDate string, groupEndDate string) (string, error) {
	query := `
		INSERT INTO groups(course_id, teacher_id, room_id, date_type, days, start_time, start_date, end_date, is_archived, name, company_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 , $11) 
		RETURNING id`

	var groupId string
	err := r.db.QueryRow(query, courseId, teacherId, roomId, dateType, pq.Array(days), lessonStartTime, groupStartDate, groupEndDate, false, name, companyId).Scan(&groupId)
	if err != nil {
		return "", err
	}
	return groupId, nil
}
func (r *GroupRepository) UpdateGroup(companyId string, id string, name string, courseId int32, teacherId string, dateType string, days []string, roomId int32, lessonStartTime string, groupStartDate string, groupEndDate string) error {
	query := `UPDATE groups SET course_id=$1, teacher_id=$2, room_id=$3, date_type=$4, days=$5, start_time=$6, start_date=$7, end_date=$8, name=$9 WHERE id=$10 and company_id=$11`
	_, err := r.db.Exec(query, courseId, teacherId, roomId, dateType, pq.Array(days), lessonStartTime, groupStartDate, groupEndDate, name, id, companyId)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupRepository) DeleteGroup(companyId string, id string) error {
	query := `UPDATE groups SET is_archived = NOT is_archived WHERE id = $1 and company_id=$2`
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupRepository) GetGroup(ctx context.Context, companyId string, page, size int32, isArchive bool) (*pb.GetGroupsResponse, error) {
	offset := (page - 1) * size

	query := `
		SELECT 
			g.id, 
			g.course_id, 
			COALESCE(c.title, 'Unknown Course') as course_title, 
			g.teacher_id,
			g.room_id, 
			COALESCE(r.title, 'Unknown Room') as room_title, 
			r.capacity, 
			g.start_date, 
			g.end_date, 
			g.is_archived, 
			g.name, 
			COUNT(gs.id) as student_count, 
			g.created_at, 
			g.days, 
			g.start_time, 
			g.date_type
		FROM groups g
		LEFT JOIN courses c ON g.course_id = c.id
		LEFT JOIN rooms r ON g.room_id = r.id
		LEFT JOIN group_students gs ON g.id = gs.group_id
		WHERE g.is_archived = $1 and g.company_id = $4
		GROUP BY g.id, c.title, r.title, r.capacity
		LIMIT $2 OFFSET $3;
	`

	countQuery := `SELECT COUNT(*) FROM groups WHERE is_archived = $1 and company_id = $2;`

	fmt.Printf("countQuery params: isArchived=%v, companyId=%v", isArchive, companyId)
	var totalCount int32
	err := r.db.QueryRow(countQuery, isArchive, companyId).Scan(&totalCount)
	if err != nil {
		log.Printf("Error counting total groups: %v (isArchived=%v, companyId=%v)", err, isArchive, companyId)
		return nil, fmt.Errorf("error counting total groups: %w", err)
	}

	totalPageCount := (totalCount + size - 1) / size

	rows, err := r.db.Query(query, isArchive, size, offset, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
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
		teacherName, err := r.userClient.GetTeacherById(ctx, group.TeacherId)
		if err != nil {
			log.Printf("Error fetching teacher by ID: %v", err)
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

	return &pb.GetGroupsResponse{
		Groups:         groups,
		TotalPageCount: totalPageCount,
	}, nil
}

func (r *GroupRepository) GetGroupById(ctx context.Context, companyId string, id, actionRole, actionId string) (*pb.GetGroupAbsResponse, error) {
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
              WHERE g.id = $1 and g.company_id=$2
              GROUP BY g.id, c.title, r.title, r.capacity , c.course_duration , c.duration_lesson , c.description , c.price`

	var group pb.GetGroupAbsResponse
	var studentCount sql.NullInt32
	var course pb.AbsCourse
	var room pb.AbsRoom

	err := r.db.QueryRow(query, id, companyId).Scan(
		&group.Id, &course.Id, &course.Name, &room.Id, &room.Name, &room.Capacity, &group.StartDate, &group.EndDate,
		&group.IsArchived, &group.Name, &studentCount, &group.CreatedAt, pq.Array(&group.Days), &group.LessonStartTime, &group.DateType, &course.CourseDuration, &course.LessonDuration, &course.Description, &course.Price, &group.TeacherId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group with id %s not found", id)
		}
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	teacherName, _ := r.userClient.GetTeacherById(ctx, group.TeacherId)
	group.TeacherName = teacherName
	group.Course = &course
	group.Room = &room
	group.StudentCount = studentCount.Int32
	return &group, nil
}
func (r *GroupRepository) GetGroupByCourseId(ctx context.Context, companyId string, courseId string) (*pb.GetGroupsByCourseResponse, error) {
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
        WHERE g.course_id = $1 and company_id=$2
    `

	rows, err := r.db.Query(query, courseId, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response pb.GetGroupsByCourseResponse
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
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

		teacherName, err := r.userClient.GetTeacherById(ctx, teacherId.String)
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
func (r *GroupRepository) GetGroupByTeacherId(companyId string, teacherId string, archived bool) (*pb.GetGroupsByTeacherResponse, error) {
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
		WHERE g.teacher_id = $1 AND g.is_archived = $2 and g.company_id=$3
	`
	rows, err := r.db.Query(query, teacherId, archived, companyId)
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
		students, err := r.GetStudentsByGroupId(companyId, group.Id)
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
func (r *GroupRepository) GetStudentsByGroupId(companyId string, groupId string) ([]*pb.AbsStudent, error) {
	query := `
		SELECT s.id, s.name, s.phone 
		FROM students s
		INNER JOIN group_students gs ON gs.student_id = s.id
		WHERE gs.group_id = $1 AND gs.condition = 'ACTIVE' and gs.company_id=$2
	`
	rows, err := r.db.Query(query, groupId, companyId)
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
func (r *GroupRepository) GetCommonInformationEducation(companyId string) (*pb.GetCommonInformationEducationResponse, error) {
	response := new(pb.GetCommonInformationEducationResponse)

	var leaveGroupCount, activeGroupCount, activeStudentCount, debtorsCount, eleminatedInTrial int32

	err := r.db.QueryRow(`SELECT COUNT(*) FROM group_students where condition='DELETE' and company_id=$1`, companyId).Scan(&leaveGroupCount)
	if err != nil {
		leaveGroupCount = 0
	}

	err = r.db.QueryRow(`SELECT COUNT(*) FROM groups where is_archived=false and company_id=$1`, companyId).Scan(&activeGroupCount)
	if err != nil {
		activeGroupCount = 0
	}

	err = r.db.QueryRow(`SELECT count(*) FROM students where condition='ACTIVE' and company_id=$1`, companyId).Scan(&activeStudentCount)
	if err != nil {
		activeStudentCount = 0
	}

	err = r.db.QueryRow(`SELECT COUNT(*) FROM students where balance < 0 and company_id=$1`, companyId).Scan(&debtorsCount)
	if err != nil {
		debtorsCount = 0
	}

	err = r.db.QueryRow(`SELECT count(*) FROM group_student_condition_history where is_eliminated_trial=true and company_id=$1`, companyId).Scan(&eleminatedInTrial)
	if err != nil {
		eleminatedInTrial = 0
	}
	response.DebtorsCount = debtorsCount
	response.LeaveGroupCount = leaveGroupCount
	response.ActiveGroupCount = activeGroupCount
	response.ActiveStudentCount = activeStudentCount
	response.EleminatedInTrial = eleminatedInTrial
	return response, nil
}
func (r *GroupRepository) GetGroupsByStudentId(companyId string, studentId string) (*pb.GetGroupsByStudentResponse, error) {
	var (
		comments []*pb.DebtorComment
		groups   []*pb.DebtorGroup
	)

	commentQuery := `
		SELECT id, comment
		FROM student_note 
		WHERE student_id = $1 and company_id=$2
	`
	commentRows, err := r.db.Query(commentQuery, studentId, companyId)
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
		WHERE gs.student_id = $1 and gs.company_id=$2
	`
	groupRows, err := r.db.Query(groupQuery, studentId, companyId)
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
func (r *GroupRepository) GetLeftAfterTrial(companyId string, from string, to string, page string, size string) (*pb.GetLeftAfterTrialPeriodResponse, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM group_student_condition_history gsch
		JOIN students ss ON gsch.student_id = ss.id
		JOIN groups g ON gsch.group_id = g.id
		WHERE gsch.is_eliminated_trial = TRUE and gsch.company_id=$3
		AND gsch.specific_date BETWEEN $1 AND $2;
	`

	var totalItemCount int
	err := r.db.QueryRow(countQuery, from, to, companyId).Scan(&totalItemCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total item count: %v", err)
	}

	query := `
		SELECT
			ss.id AS student_id,
			ss.name,
			ss.phone,
			ss.balance,
			g.id AS group_id,
			g.name AS group_name,
			gsch.return_the_money,
			gsch.created_at,
			gsch.specific_date
		FROM group_student_condition_history gsch
		JOIN students ss ON gsch.student_id = ss.id
		JOIN groups g ON gsch.group_id = g.id
		WHERE gsch.is_eliminated_trial = TRUE and gsch.company_id=$5
		AND gsch.specific_date BETWEEN $1 AND $2
		LIMIT $3 OFFSET $4;
	`

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("invalid page parameter: %v", err)
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return nil, fmt.Errorf("invalid size parameter: %v", err)
	}

	offset := (pageInt - 1) * sizeInt

	rows, err := r.db.Query(query, from, to, sizeInt, offset, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to query group_student_condition_history: %v", err)
	}
	defer rows.Close()

	var items []*pb.AbsGetLeftAfter

	for rows.Next() {
		var item pb.AbsGetLeftAfter
		var returnMoney bool
		var specificDate string
		var createdAt string

		err := rows.Scan(
			&item.StudentId,
			&item.StudentName,
			&item.StudentPhone,
			&item.StudentBalance,
			&item.GroupId,
			&item.GroupName,
			&returnMoney,
			&createdAt,
			&specificDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		item.ReturnMoney = returnMoney
		item.CreatedAt = createdAt
		item.SpecificDate = specificDate

		items = append(items, &item)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no records found")
	}

	resp := &pb.GetLeftAfterTrialPeriodResponse{
		Items:          items,
		TotalItemCount: int32(totalItemCount),
	}

	return resp, nil
}
