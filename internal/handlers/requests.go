package handlers

import (
	"encoding/json"
	"fmt"
	"go-api-tech-challenge/internal/models"
	"net/http"
)

type inputCourse struct {
	Name string `json:"name"`
}

type inputPerson struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
	Age       int    `json:"age"`
	Courses   []int  `json:"courses,omitempty"`
}

// MapTo maps a inputUser to a models.User object.
func (course inputCourse) MapTo() (models.Course, error) {
	return models.Course{
		ID:   0,
		Name: course.Name,
	}, nil
}
func (person inputPerson) MapTo() (models.Person, error) {
	return models.Person{
		ID:        0,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Type:      person.Type,
		Age:       person.Age,
		Courses:   person.Courses,
	}, nil
}

func (person inputPerson) Valid() []problem {
	var problems []problem
	validTypes := map[string]bool{
		"student":   true,
		"professor": true,
	}

	// validate FirstName is not blank
	if person.FirstName == "" {
		problems = append(problems, problem{
			Name:        "first_name",
			Description: "must not be blank",
		})
	}

	if person.LastName == "" {
		problems = append(problems, problem{
			Name:        "last_name",
			Description: "must not be blank",
		})
	}
	if !validTypes[person.Type] {
		problems = append(problems, problem{
			Name:        "type",
			Description: "must be either 'student' or 'professor'",
		})
	}
	if person.Age < 0 {
		problems = append(problems, problem{
			Name:        "age",
			Description: "must not be negative",
		})
	}
	if len(person.Courses) > 0 {
		for i, id := range person.Courses {
			if id <= 0 {
				problems = append(problems, problem{
					Name:        fmt.Sprintf("courses[%d]", i),
					Description: "course ID must be a positive integer",
				})
			}
		}
	}

	return problems
}

// Valid validates all fields of an inputUser struct.
func (course inputCourse) Valid() []problem {
	var problems []problem

	// validate FirstName is not blank
	if course.Name == "" {
		problems = append(problems, problem{
			Name:        "name",
			Description: "must not be blank",
		})
	}

	return problems
}

// problem represents an issue found during validation.
type problem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validator is an interface that defines a method for validating an object.
// It returns a slice of problems found during validation.
type Validator interface {
	Valid() (problems []problem)
}

// Mapper is a generic interface that defines a method for mapping an object to another type.
// The MapTo method returns the mapped object and an error if the mapping fails.
type Mapper[T any] interface {
	MapTo() (T, error)
}

// ValidatorMapper is a generic interface that combines Validator and Mapper interfaces.
// It requires implementing both validation and mapping methods.
type ValidatorMapper[T any] interface {
	Validator
	Mapper[T]
}

// decodeValidateBody decodes a JSON string into a ValidatorMapper, validates it, and maps it to
// the output type. If decoding, validation, or mapping fails, it returns the appropriate errors
// and problems.
func decodeValidateBody[I ValidatorMapper[O], O any](r *http.Request) (O, []problem, error) {
	var inputModel I

	// decode to JSON
	if err := json.NewDecoder(r.Body).Decode(&inputModel); err != nil {
		return *new(O), nil, fmt.Errorf("[in decodeValidateBody] decode json: %w", err)
	}

	// validate
	if problems := inputModel.Valid(); len(problems) > 0 {
		return *new(O), problems, fmt.Errorf(
			"[in decodeValidateBody] invalid %T: %d problems", inputModel, len(problems),
		)
	}

	// map to return type
	data, err := inputModel.MapTo()
	if err != nil {
		return *new(O), nil, fmt.Errorf(
			"[in decodeValidateBody] error mapping input %T to %T: %w",
			*new(I),
			*new(O),
			err,
		)
	}

	return data, nil, nil
}
