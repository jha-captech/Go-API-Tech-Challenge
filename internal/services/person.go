package services

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"jf.go.techchallenge/internal/models"
)

type PersonService struct {
	db *gorm.DB
}

func NewPersonService(db *gorm.DB) *PersonService {
	return &PersonService{
		db: db,
	}
}

func (s *PersonService) GetOneByGuid(guid string) (models.Person, error) {
	var person models.Person
	result := s.db.Table("person").First(&person, "guid = ?", guid)

	/*
		var person Person
		result := s.database.Table("person").First(&person, "guid = ?", guid)
		if result.Error != nil {
			return person, apperror.New("Person Not Found")
		}
		return person, nil
	*/
	return person, result.Error
	// if result.Error != nil {

	// }
}

func (s *PersonService) GetPersons(filters Filters) ([]models.Person, error) {

	var persons []models.Person
	tx := s.db.Table("person")

	// process query parameteres.
	for key, value := range filters {
		tx.Where(fmt.Sprintf("%s like ?", key), strings.Join([]string{"%", value, "%"}, ""))
	}

	result := tx.Find(&persons)

	return persons, result.Error
}
