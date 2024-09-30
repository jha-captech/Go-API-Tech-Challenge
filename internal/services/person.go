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

type PersonService struct {
	logger     *applog.AppLogger
	repository repository.PersonRepository
}

var validFilters = MakeFilterColumns(ValidFilters{
	"FirstName",
	"LastName",
	"Email",
})

func NewPersonService(logger *applog.AppLogger, repository repository.PersonRepository) *PersonService {
	return &PersonService{
		logger:     logger,
		repository: repository,
	}
}

func (s PersonService) GetOneByGuid(guid string) (models.Person, error) {
	person, err := s.repository.FindOne(guid)

	if err != nil {
		s.logger.Debug("Failed to query Person: ", err)
		return person, apperror.NotFound("Person: %s Not Found", guid)
	}

	return person, err
}

// Parse dont validate https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/
func parse(input models.UpdatePerson, person *models.Person) error {
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

func (s PersonService) Update(guid string, input models.UpdatePerson) (models.Person, error) {
	person, err := s.GetOneByGuid(guid)
	if err != nil {
		return person, err
	}

	err = parse(input, &person)

	if err != nil {
		return person, err
	}

	err = s.repository.Save(&person)

	if err != nil {
		// todo log and return different error
	}

	return person, err
}

func (s PersonService) Delete(guid string) error {
	person, err := s.GetOneByGuid(guid)
	if err != nil {
		return err
	}
	err = s.repository.Delete(&person)

	if err != nil {
		// todo log and return database error
	}
	return err
}

func (s PersonService) Create(input models.UpdatePerson) (models.Person, error) {
	newPerson := models.Person{
		Guid: uuid.NewString(),
	}

	err := parse(input, &newPerson)

	if err != nil {
		return newPerson, err
	}

	err = s.repository.Save(&newPerson)

	if err != nil {
		// todo log and return different error..
	}
	return newPerson, err
}

func (s PersonService) GetAll(urlParams url.Values) ([]models.Person, error) {

	filters, err := ParseURLFilters(urlParams, validFilters)

	if err != nil {
		return nil, err
	}

	persons, err := s.repository.FindAll(filters)

	if err != nil {
		// todo handle database error..
	}

	return persons, nil
}
