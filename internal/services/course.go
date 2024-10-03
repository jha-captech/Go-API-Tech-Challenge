package services

import (
	"context"
	"database/sql"
	"fmt"
	"go-api-tech-challenge/internal/models"
)

type CourseService struct {
	database *sql.DB
}

func NewCourseService(db *sql.DB) *CourseService {
	return &CourseService{
		database: db,
	}
}

func (s *CourseService) ListCourses(ctx context.Context) ([]models.Course, error) {

	query := `SELECT * FROM course 
	ORDER BY id asc`
	rows, err := s.database.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err = rows.Scan(&course.ID, &course.Name)
		if err != nil {
			return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to scan course from row: %w", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to scan courses: %w", err)
	}

	return courses, nil

}

func (s *CourseService) GetCourseByID(ctx context.Context, id int) (models.Course, error) {
	var course models.Course
	query := "SELECT id, name FROM course WHERE id = $1"

	err := s.database.QueryRowContext(ctx, query, id).Scan(&course.ID, &course.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Course{}, fmt.Errorf("[in services.GetCourseByIDByID] no course found with id: %d", id)
		}
		return models.Course{}, fmt.Errorf("[in services.GetCourseByIDByID] failed to retrieve course: %w", err)
	}

	return course, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, courseID int, newName string) (models.Course, error) {
	query := `UPDATE course SET name = $1 WHERE id = $2`
	result, err := s.database.ExecContext(ctx, query, newName, courseID)
	if err != nil {

		return models.Course{}, fmt.Errorf("[in services.UpdateCourse] failed to update course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.UpdateCourse] failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.Course{}, fmt.Errorf("[in services.UpdateCourse] no course found with id: %d", courseID)

	}

	return models.Course{
		ID:   courseID,
		Name: newName,
	}, nil
}

func (s *CourseService) CreateCourse(ctx context.Context, courseName string) (models.Course, error) {
	query := `INSERT INTO course (name) VALUES ($1) RETURNING id`
	var newID int

	err := s.database.QueryRowContext(ctx, query, courseName).Scan(&newID)
	if err != nil {
		return models.Course{}, fmt.Errorf("failed to create course: %w", err)
	}
	return models.Course{ID: newID, Name: courseName}, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, courseID int) error {
	query := `DELETE FROM course WHERE id = $1`

	result, err := s.database.ExecContext(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("[in services.ListCourses] failed to delete course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[in services.ListCourses] failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("[in services.ListCourses] no course found with id %d", courseID)
	}

	return nil
}
