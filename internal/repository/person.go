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

	Save(person *models.Person, inputCourses []models.Course) error

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
// In order for courses to be processed correctly, the persons courses must be "hydrated"
func (s PersonRepositoryImpl) Save(person *models.Person, inputCourses []models.Course) error {

	// var coursesToSave []models.PersonCourse
	var coursesToDelete []models.PersonCourse

	// Tracks courses state requested by user
	courseMap := map[uint]interface{}{}
	for _, value := range inputCourses {
		courseMap[value.ID] = nil
	}

	// If any person course is not present in courses Map, it should be deleted.
	for _, course := range person.Courses {
		_, present := courseMap[course.ID]
		if !present {
			coursesToDelete = append(coursesToDelete, models.PersonCourse{
				PersonID: person.ID,
				CourseID: course.ID,
			})
		}
	}

	// // // Iterate both collections so all courses are represented.
	// for _, course := range append(person.Courses, inputCourses...) {
	// 	courseMap[course.ID] = false
	// }

	// // Update courses to be saved from input course.
	// for _, inputCourse := range inputCourses {
	// 	_, present := courseMap[inputCourse.ID]
	// 	if present {
	// 		courseMap[inputCourse.ID] = true
	// 		coursesToSave = append(coursesToSave, models.PersonCourse{
	// 			PersonID: person.ID,
	// 			CourseID: inputCourse.ID,
	// 		})
	// 	}
	// }

	// // Courses to delete
	// for key, value := range courseMap {
	// 	if !value {
	// 		coursesToDelete = append(coursesToDelete, models.PersonCourse{
	// 			PersonID: person.ID,
	// 			CourseID: key,
	// 		})
	// 	}
	// }

	person.Courses = inputCourses

	return logDBErr(s.logger, s.db.Transaction(func(tx *gorm.DB) error {

		var coursesError error
		if len(coursesToDelete) > 0 {
			coursesError = tx.Delete(&coursesToDelete).Error
		}

		if coursesError != nil {
			return coursesError
		}

		// if len(coursesToSave) > 0 {
		// 	coursesError = tx.Save(&coursesToSave).Error
		// }

		// if coursesError != nil {
		// 	return coursesError
		// }

		return tx.Save(person).Error
	}), "Failed to Save Person!")
}

func (s PersonRepositoryImpl) Delete(person *models.Person) error {

	// Handle deleting the courses the person is enrolled in as well as the person.
	return logDBErr(s.logger, s.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Delete(&models.PersonCourse{}, "person_id = ?", person.ID)
		if err := logDBErr(s.logger, result.Error, "Failed to delete person_course record"); err != nil {
			return err
		}

		return tx.Delete(person).Error
	}), "Failed to delete person record")
}
