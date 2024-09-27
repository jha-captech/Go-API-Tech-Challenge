package services

import (
	"fmt"

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

// func (s *PersonSe)

func (s *PersonService) GetPersons(filters map[string]string) ([]models.Person, error) {

	var persons []models.Person
	tx := s.db.Table("person")

	// process query parameteres.
	for key, value := range filters {
		tx.Where(fmt.Sprintf("%s <> ?", key), value)
	}

	result := tx.Find(&persons)

	return persons, result.Error
}
