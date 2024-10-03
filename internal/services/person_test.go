package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"go-api-tech-challenge/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type personTestSuite struct {
	suite.Suite
	service *PersonService
	dbMock  sqlmock.Sqlmock
}

func TestPersonTestSuite(t *testing.T) {
	suite.Run(t, new(personTestSuite))
}

func (s *personTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	s.dbMock = mock
	s.service = NewPersonService(db)
}

func (s *personTestSuite) TearDownSuite() {
	err := s.dbMock.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

func (s *personTestSuite) TestListPersons() {
	t := s.T()

	persons := []models.Person{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Type:      "student",
			Age:       25,
			Courses:   []int{1, 2},
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Type:      "professor",
			Age:       45,
			Courses:   []int{3},
		},
	}

	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		expectedReturn []models.Person
		expectedError  error
	}{
		"Return slice of persons": {
			mockReturn: sqlmock.NewRows([]string{"person_id", "first_name", "last_name", "type", "age", "course_ids"}).
				AddRow(1, "John", "Doe", "student", 25, pq.Array([]int64{1, 2})).
				AddRow(2, "Jane", "Smith", "professor", 45, pq.Array([]int64{3})),
			mockReturnErr:  nil,
			expectedReturn: persons,
			expectedError:  nil,
		},
		"Error getting persons": {
			mockReturn:     sqlmock.NewRows([]string{}),
			mockReturnErr:  errors.New("test error"),
			expectedReturn: []models.Person{},
			expectedError:  fmt.Errorf("[in services.ListPersons] failed to get persons: %w", errors.New("test error")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `SELECT p.id as person_id, p.first_name, p.last_name, p.type, p.age, COALESCE(Array_AGG(pc.course_id), '{}') as course_ids
				FROM person p
				LEFT JOIN person_course pc ON p.id = pc.person_id
				GROUP BY p.id, p.first_name, p.last_name, p.type, p.age
				ORDER BY person_id asc`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.ListPersons(context.Background())

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedReturn, actualReturn)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func (s *personTestSuite) TestGetPersonByName() {
	t := s.T()

	person := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	testCases := map[string]struct {
		name           string
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		expectedReturn models.Person
		expectedError  error
	}{
		"person found": {
			name: "Doe",
			mockReturn: sqlmock.NewRows([]string{"person_id", "first_name", "last_name", "type", "age", "course_ids"}).
				AddRow(1, "John", "Doe", "student", 25, pq.Array([]int64{1, 2})),
			mockReturnErr:  nil,
			expectedReturn: person,
			expectedError:  nil,
		},
		"person not found": {
			name:           "Smith",
			mockReturn:     sqlmock.NewRows([]string{}), // No rows returned
			mockReturnErr:  sql.ErrNoRows,               // Simulate no rows found
			expectedReturn: models.Person{},
			expectedError:  fmt.Errorf("[in services.GetPersonByName] no person found with name: %s", "Smith"),
		},
		"query error": {
			name:           "Doe",
			mockReturn:     sqlmock.NewRows([]string{}), // Return empty rows but simulate an error
			mockReturnErr:  errors.New("test error"),
			expectedReturn: models.Person{},
			expectedError:  fmt.Errorf("[in services.GetPersonByName] failed to retrieve person: %w", errors.New("test error")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `
				SELECT p.id as person_id, 
				p.first_name, 
				p.last_name,
				p.type,
				p.age,
				COALESCE(Array_AGG(pc.course_id), '{}') as course_ids
				FROM person p
				LEFT JOIN person_course pc ON p.id = pc.person_id
				WHERE LOWER(p.last_name) = LOWER($1)
				GROUP BY p.id, p.first_name, p.last_name, p.type, p.age`

			query := s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WithArgs(tc.name)

			if tc.mockReturnErr != nil {
				query.WillReturnError(tc.mockReturnErr)
			} else {
				query.WillReturnRows(tc.mockReturn)
			}

			actualReturn, err := s.service.GetPersonByName(context.Background(), tc.name)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedReturn, actualReturn)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *personTestSuite) TestUpdatePerson() {
	t := s.T()

	lastName := "Doe"
	personIn := models.Person{
		FirstName: "John",
		LastName:  "Smith",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	personOut := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	s.dbMock.ExpectBegin()

	updatePersonQuery := `
		UPDATE person
		SET first_name = $1,
			last_name = $2,
			type = $3,
			age = $4
		WHERE LOWER(last_name) = LOWER($5)
		RETURNING id, first_name, last_name, type, age;
	`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(updatePersonQuery)).
		WithArgs(personIn.FirstName, personIn.LastName, personIn.Type, personIn.Age, lastName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
			AddRow(personOut.ID, personOut.FirstName, personOut.LastName, personOut.Type, personOut.Age))

	deleteCoursesQuery := `DELETE FROM person_course WHERE person_id = $1`
	s.dbMock.ExpectExec(regexp.QuoteMeta(deleteCoursesQuery)).
		WithArgs(personOut.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertCoursesQuery := `INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)`
	for _, courseID := range personIn.Courses {
		s.dbMock.ExpectExec(regexp.QuoteMeta(insertCoursesQuery)).
			WithArgs(personOut.ID, courseID).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	selectCoursesQuery := `SELECT course_id FROM person_course WHERE person_id = $1`
	s.dbMock.ExpectQuery(regexp.QuoteMeta(selectCoursesQuery)).
		WithArgs(personOut.ID).
		WillReturnRows(sqlmock.NewRows([]string{"course_id"}).
			AddRow(1).
			AddRow(2))

	s.dbMock.ExpectCommit()

	result, err := s.service.UpdatePerson(context.Background(), lastName, personIn)

	assert.NoError(t, err)
	assert.Equal(t, personOut, result)

	err = s.dbMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func (s *personTestSuite) TestCreatePerson() {
	t := s.T()

	personIn := models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	personOut := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []int{1, 2},
	}

	t.Run("person created successfully", func(t *testing.T) {
		s.dbMock.ExpectBegin()

		insertPersonQuery := `
			INSERT INTO person (first_name, last_name, type, age) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id, first_name, last_name, type, age
		`
		s.dbMock.ExpectQuery(regexp.QuoteMeta(insertPersonQuery)).
			WithArgs(personIn.FirstName, personIn.LastName, personIn.Type, personIn.Age).
			WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
				AddRow(personOut.ID, personOut.FirstName, personOut.LastName, personOut.Type, personOut.Age))

		insertCourseQuery := `INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)`
		for _, courseID := range personIn.Courses {
			s.dbMock.ExpectExec(regexp.QuoteMeta(insertCourseQuery)).
				WithArgs(personOut.ID, courseID).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		selectCoursesQuery := `SELECT course_id FROM person_course WHERE person_id = $1`
		s.dbMock.ExpectQuery(regexp.QuoteMeta(selectCoursesQuery)).
			WithArgs(personOut.ID).
			WillReturnRows(sqlmock.NewRows([]string{"course_id"}).
				AddRow(1).
				AddRow(2))

		s.dbMock.ExpectCommit()

		actualReturn, err := s.service.CreatePerson(context.Background(), personIn)

		assert.NoError(t, err)
		assert.Equal(t, personOut, actualReturn)

		err = s.dbMock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func (s *personTestSuite) TestDeletePerson() {
	t := s.T()

	lastName := "Doe"

	testCases := map[string]struct {
		mockBeginErr         error
		mockDeleteCoursesErr error
		mockDeletePersonErr  error
		mockRowsAffected     int64
		mockCommitErr        error
		expectedError        error
	}{
		"successful deletion": {
			mockBeginErr:         nil,
			mockDeleteCoursesErr: nil,
			mockDeletePersonErr:  nil,
			mockRowsAffected:     1,
			mockCommitErr:        nil,
			expectedError:        nil,
		},
		"begin transaction error": {
			mockBeginErr:         errors.New("transaction begin error"),
			mockDeleteCoursesErr: nil,
			mockDeletePersonErr:  nil,
			mockRowsAffected:     0,
			mockCommitErr:        nil,
			expectedError:        fmt.Errorf("failed to begin transaction: %w", errors.New("transaction begin error")),
		},
		"error deleting courses": {
			mockBeginErr:         nil,
			mockDeleteCoursesErr: errors.New("delete courses error"),
			mockDeletePersonErr:  nil,
			mockRowsAffected:     0,
			mockCommitErr:        nil,
			expectedError:        fmt.Errorf("failed to delete courses: %w", errors.New("delete courses error")),
		},
		"error deleting person": {
			mockBeginErr:         nil,
			mockDeleteCoursesErr: nil,
			mockDeletePersonErr:  errors.New("delete person error"),
			mockRowsAffected:     0,
			mockCommitErr:        nil,
			expectedError:        fmt.Errorf("failed to delete person: %w", errors.New("delete person error")),
		},
		"no person found": {
			mockBeginErr:         nil,
			mockDeleteCoursesErr: nil,
			mockDeletePersonErr:  nil,
			mockRowsAffected:     0,
			mockCommitErr:        nil,
			expectedError:        fmt.Errorf("no person found with last name: %s", lastName),
		},
		"commit transaction error": {
			mockBeginErr:         nil,
			mockDeleteCoursesErr: nil,
			mockDeletePersonErr:  nil,
			mockRowsAffected:     1,
			mockCommitErr:        errors.New("commit transaction error"),
			expectedError:        fmt.Errorf("failed to commit transaction: %w", errors.New("commit transaction error")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.mockBeginErr != nil {
				s.dbMock.ExpectBegin().WillReturnError(tc.mockBeginErr)
			} else {
				s.dbMock.ExpectBegin()
			}

			if tc.mockBeginErr == nil {
				deleteCoursesQuery := `
					DELETE FROM person_course WHERE person_id IN 
					(SELECT id FROM person WHERE LOWER(last_name) = LOWER($1))
				`
				if tc.mockDeleteCoursesErr != nil {
					s.dbMock.ExpectExec(regexp.QuoteMeta(deleteCoursesQuery)).
						WithArgs(lastName).
						WillReturnError(tc.mockDeleteCoursesErr)
				} else {
					s.dbMock.ExpectExec(regexp.QuoteMeta(deleteCoursesQuery)).
						WithArgs(lastName).
						WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate course deletion
				}

				if tc.mockDeleteCoursesErr == nil {
					deletePersonQuery := `DELETE FROM person WHERE LOWER(last_name) = LOWER($1)`
					if tc.mockDeletePersonErr != nil {
						s.dbMock.ExpectExec(regexp.QuoteMeta(deletePersonQuery)).
							WithArgs(lastName).
							WillReturnError(tc.mockDeletePersonErr)
					} else {
						s.dbMock.ExpectExec(regexp.QuoteMeta(deletePersonQuery)).
							WithArgs(lastName).
							WillReturnResult(sqlmock.NewResult(1, tc.mockRowsAffected)) // Mock rows affected
					}

					if tc.mockDeletePersonErr == nil && tc.mockRowsAffected > 0 {
						// Mock commit
						if tc.mockCommitErr != nil {
							s.dbMock.ExpectCommit().WillReturnError(tc.mockCommitErr)
						} else {
							s.dbMock.ExpectCommit()
						}
					}
				}
			}

			err := s.service.DeletePerson(context.Background(), lastName)

			assert.Equal(t, tc.expectedError, err)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
