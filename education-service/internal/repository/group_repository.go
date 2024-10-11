package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r GroupRepository) CreateGroup(name string, courseId int32, teacherId string, dateType string, days []string, roomId int32, lessonStartTime string, groupStartDate string, groupEndDate string) error {
	query := `INSERT INTO groups(course_id, teacher_id, room_id, date_type, days, start_time, start_date, end_date, is_archived,name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,  $10)`
	_, err := r.db.Exec(query, courseId, teacherId, roomId, dateType, pq.Array(days), lessonStartTime, groupStartDate, groupEndDate, false, name)
	if err != nil {
		return err
	}
	return nil
}

func (r GroupRepository) UpdateGroup(id string, name string, courseId int32, teacherId string, dateType string, days []string, roomId int32, lessonStartTime string, groupStartDate string, groupEndDate string) error {
	query := `UPDATE groups SET course_id=$1, teacher_id=$2, room_id=$3, date_type=$4, days=$5, start_time=$6, start_date=$7, end_date=$8, name=$9 WHERE id=$10`
	_, err := r.db.Exec(query, courseId, teacherId, roomId, dateType, pq.Array(days), lessonStartTime, groupStartDate, groupEndDate, name, id)
	if err != nil {
		return err
	}
	return nil
}

func (r GroupRepository) DeleteGroup(id string) error {
	query := `DELETE FROM groups WHERE id=$1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r GroupRepository) GetGroup() (*pb.GetGroupsResponse, error) {
	query := `SELECT g.id, g.course_id, COALESCE(c.title, 'Unknown Course') as course_title, g.teacher_id, 'something' as teacher_name, g.room_id, COALESCE(r.title, 'Unknown Room') as room_title, g.date_type, g.start_time, g.start_date, g.end_date, g.is_archived, g.name, COUNT(gs.id) as student_count, g.created_at , g.start_time
              FROM groups g
              LEFT JOIN courses c ON g.course_id = c.id
              LEFT JOIN rooms r ON g.room_id = r.id
              LEFT JOIN group_students gs ON g.id = gs.group_id
              GROUP BY g.id, c.title, r.title`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*pb.GetGroupAbsResponse
	for rows.Next() {
		var group pb.GetGroupAbsResponse
		var dateType string
		var startTime string
		var courseTitle, roomTitle string
		var studentCount int32

		err := rows.Scan(&group.Id, &group.Course.Id, &courseTitle, &group.TeacherName, &group.Room.Id, &roomTitle, dateType, &group.TimeDays, &group.StartDate, &group.EndDate, &group.IsArchived, &group.Name, &studentCount, &group.CreatedAt, &startTime)
		if err != nil {
			return nil, err
		}
		group.Course.Name = courseTitle
		group.Room.Name = roomTitle
		group.StudentCount = studentCount
		group.TimeDays = dateType + " " + startTime

		groups = append(groups, &group)
	}

	return &pb.GetGroupsResponse{Groups: groups}, nil
}

func (r GroupRepository) GetGroupById(id string) (*pb.GetGroupAbsResponse, error) {
	query := `SELECT g.id, g.course_id, COALESCE(c.title, 'Unknown Course') as course_title, g.teacher_id, 'something' as teacher_name, g.room_id, COALESCE(r.title, 'Unknown Room') as room_title, g.date_type, g.start_time, g.start_date, g.end_date, g.is_archived, g.name, COUNT(gs.id) as student_count, g.created_at
              FROM groups g
              LEFT JOIN courses c ON g.course_id = c.id
              LEFT JOIN rooms r ON g.room_id = r.id
              LEFT JOIN group_students gs ON g.id = gs.group_id
              WHERE g.id = $1
              GROUP BY g.id, c.title, r.title`

	row := r.db.QueryRow(query, id)

	var group pb.GetGroupAbsResponse
	var dateType, startTime, courseTitle, roomTitle string
	var studentCount int32

	err := row.Scan(&group.Id, &group.Course.Id, &courseTitle, &group.TeacherName, &group.Room.Id, &roomTitle, &dateType, &startTime, &group.StartDate, &group.EndDate, &group.IsArchived, &group.Name, &studentCount, &group.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group with id %s not found", id)
		}
		return nil, err
	}
	group.Course.Name = courseTitle
	group.Room.Name = roomTitle
	group.StudentCount = studentCount
	group.TimeDays = dateType + " " + startTime

	return &group, nil
}
