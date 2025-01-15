package repository

import (
	"database/sql"
	"education-service/proto/pb"
	"fmt"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(companyId, title, description string, durationLesson, courseDuration int32, price float64) error {
	query := "INSERT INTO courses (title, duration_lesson, course_duration, price, description , company_id) VALUES ($1 , $2 , $3 , $4 , $5 , $6)"
	_, err := r.db.Exec(query, title, durationLesson, courseDuration, price, description, companyId)
	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

func (r *CourseRepository) UpdateCourse(companyId, title, description, id string, durationLesson, courseDuration int32, price float64) error {
	query := "UPDATE courses SET title=$1, duration_lesson=$2, course_duration=$3, price=$4, description=$5  WHERE id = $6 and company_id=$7"
	_, err := r.db.Exec(query, title, durationLesson, courseDuration, price, description, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}
	return nil
}

func (r *CourseRepository) DeleteCourse(companyId, id string) error {
	query := "DELETE FROM courses WHERE id = $1 and company_id=$2"
	_, err := r.db.Exec(query, id, companyId)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}

func (r *CourseRepository) GetCourse(companyId string) (*pb.GetUpdateCourseAbs, error) {
	query := "SELECT id, title, duration_lesson, course_duration, price, description FROM courses where company_id=$1"
	rows, err := r.db.Query(query, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result pb.GetUpdateCourseAbs
	var courses []*pb.AbsCourse
	for rows.Next() {
		var res pb.AbsCourse
		err := rows.Scan(&res.Id, &res.Name, &res.LessonDuration, &res.CourseDuration, &res.Price, &res.Description)
		if err != nil {
			return nil, err
		}
		courses = append(courses, &res)
	}
	result.Courses = courses
	return &result, nil
}

func (r *CourseRepository) GetCourseById(companyId, id string) (*pb.GetCourseByIdResponse, error) {
	query := "SELECT id, title, duration_lesson, course_duration, price, description FROM courses WHERE id = $1 and company_id=$2"
	var response pb.GetCourseByIdResponse
	err := r.db.QueryRow(query, id, companyId).Scan(&response.Id, &response.Name, &response.LessonDuration, &response.CourseDuration, &response.Price, &response.Description)
	if err != nil {
		return nil, err
	}
	response.StudentCount = 5
	return &response, nil
}
