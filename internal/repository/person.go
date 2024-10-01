package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
)

type Person interface {
	FindAll(filters Filters) ([]models.Person, error)

	FindOne(guid string) (models.Person, error)

	Save(person *models.Person) error

	Delete(person *models.Person) error
}

func NewPerson(db *gorm.DB, logger *applog.AppLogger) Person {
	return &PersonRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

type PersonRepositoryImpl struct {
	db     *gorm.DB
	logger *applog.AppLogger
}

func (s PersonRepositoryImpl) FindAll(filters Filters) ([]models.Person, error) {
	var persons []models.Person
	tx := s.db.Table("person")

	// process query parameteres.
	for key, value := range filters {
		tx.Where(fmt.Sprintf("%s like ?", key), strings.Join([]string{"%", value, "%"}, ""))
	}

	result := tx.Find(&persons)
	return persons, logDBErr(s.logger, result.Error, "Failed to Query Person Table")
}

func (s PersonRepositoryImpl) FindOne(guid string) (models.Person, error) {
	var person models.Person

	result := s.db.Table("person").Preload("Courses").Find(&person, "guid = ?", guid)

	if result.RowsAffected == 0 {
		return person, apperror.NotFound("Person: %s Not Found", guid)
	}

	return person, logDBErr(s.logger, result.Error, "Failed to Query Person Table")
}

// Used for both update and insert.
func (s PersonRepositoryImpl) Save(person *models.Person) error {
	return s.db.Save(person).Error
}

func (s PersonRepositoryImpl) Delete(person *models.Person) error {

	// Handle deleting the courses the person is enrolled in as well as the person.
	return s.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Delete(&models.PersonCourse{}, "person_id = ?", person.ID)
		if err := logDBErr(s.logger, result.Error, "Failed to delete person_course record"); err != nil {
			return err
		}

		return logDBErr(s.logger, tx.Delete(person).Error, "Failed to delete person record")
	})
}
