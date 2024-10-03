package services

import (
	"net/url"
	"strings"

	"github.com/google/uuid"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/repository"
)

type Course struct {
	logger     *applog.AppLogger
	repository repository.Course
}

func NewCourse(logger *applog.AppLogger, repository repository.Course) *Course {
	return &Course{
		logger:     logger,
		repository: repository,
	}
}

func (s Course) GetOneByGuid(guid string) (models.Course, error) {
	return s.repository.FindOne(guid)
}

func (s Course) Update(guid string, input models.CourseInput) (models.Course, error) {
	course, err := s.GetOneByGuid(guid)
	if err != nil {
		return course, err
	}

	err = s.parse(input, &course)

	if err != nil {
		return course, err
	}

	err = s.repository.Save(&course)
	return course, err
}

func (s Course) Delete(guid string) error {
	course, err := s.GetOneByGuid(guid)
	if err != nil {
		return err
	}

	return s.repository.Delete(&course)
}

func (s Course) Create(input models.CourseInput) (models.Course, error) {
	newCourse := models.Course{}

	err := s.parse(input, &newCourse)

	if err != nil {
		return newCourse, err
	}
	newCourse.Guid = uuid.NewString()
	err = s.repository.Save(&newCourse)
	return newCourse, err
}

var courseFilters = MakeFilterColumns(ValidFilters{
	"Name",
})

func (s Course) GetAll(urlParams url.Values) ([]models.Course, error) {
	filters, err := ParseURLFilters(urlParams, courseFilters)

	if err != nil {
		return nil, err
	}

	return s.repository.FindAll(filters)
}

func (s Course) parse(input models.CourseInput, course *models.Course) error {
	var errors []error

	if strings.Trim(input.Name, " ") == "" {
		errors = append(errors, apperror.BadRequest("Name must not be blank"))
	}

	if len(errors) > 0 {
		return apperror.Of(errors)
	}

	course.Name = input.Name

	return nil
}
