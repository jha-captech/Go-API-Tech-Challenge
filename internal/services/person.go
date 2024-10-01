package services

import (
	"net/mail"
	"net/url"
	"slices"
	"strings"

	"github.com/google/uuid"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/repository"
)

type Person struct {
	logger     *applog.AppLogger
	repository repository.Person
}

func NewPerson(logger *applog.AppLogger, repository repository.Person) *Person {
	return &Person{
		logger:     logger,
		repository: repository,
	}
}

// Parse dont validate https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/
func (s Person) parse(input models.UpdatePerson, person *models.Person) error {
	var errors []error

	if strings.Trim(input.FirstName, " ") == "" {
		errors = append(errors, apperror.BadRequest("First Name must not be blank"))
	}

	if strings.Trim(input.LastName, " ") == "" {
		errors = append(errors, apperror.BadRequest("Last Name must not be blank"))
	}

	if strings.Trim(input.FirstName, " ") == "" {
		errors = append(errors, apperror.BadRequest("First Name must not be blank"))
	}

	// Validate email address format
	_, emailErr := mail.ParseAddress(input.Email)
	if emailErr != nil {
		errors = append(errors, apperror.BadRequest("Email must be a valid email adress"))
	}

	if input.Age < 10 {
		errors = append(errors, apperror.BadRequest("Must be at least 10 years old to enrol."))
	}

	if slices.Contains([]models.PersonType{models.Professor, models.Strudent}, input.Type) {
		errors = append(errors, apperror.BadRequest("Invalid Person type, must be either 'professor' or 'student'"))
	}

	if len(errors) > 0 {
		return apperror.Of(errors...)
	}

	person.FirstName = input.FirstName
	person.LastName = input.LastName
	person.Email = input.Email
	person.Type = string(input.Type)
	person.Age = input.Age

	return nil
}

func (s Person) GetOneByGuid(guid string) (models.Person, error) {
	return s.repository.FindOne(guid)
}

func (s Person) Update(guid string, input models.UpdatePerson) (models.Person, error) {
	person, err := s.GetOneByGuid(guid)
	if err != nil {
		return person, err
	}

	err = s.parse(input, &person)

	if err != nil {
		return person, err
	}

	err = s.repository.Save(&person)
	return person, err
}

func (s Person) Delete(guid string) error {
	person, err := s.GetOneByGuid(guid)
	if err != nil {
		return err
	}
	err = s.repository.Delete(&person)

	return err
}

func (s Person) Create(input models.UpdatePerson) (models.Person, error) {
	newPerson := models.Person{}

	err := s.parse(input, &newPerson)

	if err != nil {
		return newPerson, err
	}
	newPerson.Guid = uuid.NewString()
	err = s.repository.Save(&newPerson)
	return newPerson, err
}

var personFilters = MakeFilterColumns(ValidFilters{
	"FirstName",
	"LastName",
	"Email",
})

func (s Person) GetAll(urlParams url.Values) ([]models.Person, error) {

	filters, err := ParseURLFilters(urlParams, personFilters)

	if err != nil {
		return nil, err
	}

	return s.repository.FindAll(filters)
}
