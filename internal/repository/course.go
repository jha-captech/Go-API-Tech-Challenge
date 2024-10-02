package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
)

type Course interface {
	FindAll(filters Filters) ([]models.Course, error)

	FindOne(guid string) (models.Course, error)

	Save(course *models.Course) error

	Delete(course *models.Course) error
}

func NewCourse(db *gorm.DB, logger *applog.AppLogger) Course {
	return CourseRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

type CourseRepositoryImpl struct {
	db     *gorm.DB
	logger *applog.AppLogger
}

func (s CourseRepositoryImpl) FindAll(filters Filters) ([]models.Course, error) {
	var courses []models.Course
	tx := s.db.Table("course")

	for key, value := range filters {
		tx.Where(fmt.Sprintf("%s like ?", key), strings.Join([]string{"%", value, "%"}, ""))
	}

	result := tx.Find(&courses)
	return courses, LogDBErr(s.logger, result.Error, "Failed to query courses table")
}

func (s CourseRepositoryImpl) FindOne(guid string) (models.Course, error) {
	var course models.Course

	result := s.db.Table("course").Find(&course, "guid = ?", guid)

	if result.RowsAffected == 0 {
		return course, apperror.NotFound("Course: %s Not Found", guid)
	}

	return course, LogDBErr(s.logger, result.Error, "Failed to Query Course Table")
}

func (s CourseRepositoryImpl) Save(course *models.Course) error {
	return LogDBErr(s.logger, s.db.Save(course).Error, "Failed to Save Course")
}

func (s CourseRepositoryImpl) Delete(course *models.Course) error {
	// Handle deleting the courses
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.PersonCourse{}, "course_id = ?", course.ID)

		if err := LogDBErr(s.logger, result.Error, "Failed to delete person_course for course record"); err != nil {
			return err
		}

		return LogDBErr(s.logger, tx.Delete(course).Error, "Failed to delete course record")
	})
}
