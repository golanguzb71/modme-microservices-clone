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

func (r *CourseRepository) CreateCourse(title, description string, durationLesson, courseDuration int32, price float64) error {
	query := "INSERT INTO courses (title, duration_lesson, course_duration, price, description) VALUES ($1 , $2 , $3 , $4 , $5)"
	_, err := r.db.Exec(query, title, durationLesson, courseDuration, price, description)
	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

func (r *CourseRepository) UpdateCourse(title, description, id string, durationLesson, courseDuration int32, price float64) error {
	query := "UPDATE courses SET title=$1, duration_lesson=$2, course_duration=$3, price=$4, description=$5  WHERE id = $6"
	_, err := r.db.Exec(query, title, durationLesson, courseDuration, price, description, id)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}
	return nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	query := "DELETE FROM courses WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}

func (r *CourseRepository) GetCourse() (*pb.GetUpdateCourseAbs, error) {
	query := "SELECT id, title, duration_lesson, course_duration, price, description FROM courses"
	rows, err := r.db.Query(query)
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

func (r *CourseRepository) GetCourseById(id string) (*pb.GetCourseByIdResponse, error) {
	query := "SELECT id, title, duration_lesson, course_duration, price, description FROM courses WHERE id = $1"
	var response pb.GetCourseByIdResponse
	err := r.db.QueryRow(query, id).Scan(&response.Id, &response.Name, &response.LessonDuration, &response.CourseDuration, &response.Price, &response.Description)
	if err != nil {
		return nil, err
	}
	response.StudentCount = 5
	return &response, nil
}
