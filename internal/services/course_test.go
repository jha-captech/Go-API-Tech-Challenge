package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"go-api-tech-challenge/internal/models"
	"go-api-tech-challenge/internal/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuit struct {
	suite.Suite
	service *CourseService
	dbMock  sqlmock.Sqlmock
}

func TestTestSuit(t *testing.T) {
	suite.Run(t, new(testSuit))
}

func (s *testSuit) SetupSuite() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	s.dbMock = mock
	s.service = NewCourseService(db)
}

func (s *testSuit) TearDownSuite() {
	err := s.dbMock.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

func (s *testSuit) TestListCourses() {
	t := s.T()

	courses := []models.Course{
		{ID: 1, Name: "Databases"},
		{ID: 2, Name: "Operating Systems"},
	}

	testCases := map[string]struct {
		mockReturn     *sqlmock.Rows
		mockReturnErr  error
		expectedReturn []models.Course
		expectedError  error
	}{
		"Return slice of courses": {
			mockReturn:     testutil.MustStructsToRows(courses),
			mockReturnErr:  nil,
			expectedReturn: courses,
			expectedError:  nil,
		},
		"Error getting courses": {
			mockReturn:     &sqlmock.Rows{},
			mockReturnErr:  errors.New("test"),
			expectedReturn: []models.Course{},
			expectedError:  fmt.Errorf("[in services.ListCourses] failed to get courses: %w", errors.New("test")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `SELECT * FROM course`
			s.dbMock.
				ExpectQuery(regexp.QuoteMeta(exp)).
				WillReturnRows(tc.mockReturn).
				WillReturnError(tc.mockReturnErr)

			actualReturn, err := s.service.ListCourses(context.Background())

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedReturn, actualReturn)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestUpdateCourse() {
	t := s.T()

	courseIn := models.Course{ID: 1, Name: "Databases"}
	courseOut := models.Course{ID: 1, Name: "Advanced Databases"}

	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockReturn     driver.Result
		mockReturnErr  error
		inputID        int
		inputCourse    models.Course
		expectedReturn models.Course
		expectedError  error
	}{
		"course updated by ID": {
			mockInputArgs:  []driver.Value{courseOut.Name, int(courseOut.ID)},
			mockReturn:     sqlmock.NewResult(1, 1),
			mockReturnErr:  nil,
			inputID:        int(courseIn.ID),
			inputCourse:    courseOut,
			expectedReturn: courseOut,
			expectedError:  nil,
		},
		"Error updating course": {
			mockInputArgs:  []driver.Value{courseIn.Name, 5},
			mockReturn:     nil,
			mockReturnErr:  errors.New("test"),
			inputID:        5,
			inputCourse:    courseIn,
			expectedReturn: models.Course{},
			expectedError:  fmt.Errorf("[in services.UpdateCourse] failed to update course: %w", errors.New("test")),
		},
		"no rows affected": {
			mockInputArgs:  []driver.Value{courseIn.Name, 88},
			mockReturn:     sqlmock.NewResult(0, 0),
			mockReturnErr:  nil,
			inputID:        88,
			inputCourse:    courseIn,
			expectedReturn: models.Course{},
			expectedError:  fmt.Errorf("[in services.UpdateCourse] no course found with id: %d", 88),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			exp := `UPDATE course 
			SET name = $1 
			WHERE id = $2`
			mock := s.dbMock.ExpectExec(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...)

			if tc.mockReturnErr != nil {
				mock.WillReturnError(tc.mockReturnErr)
			} else {
				mock.WillReturnResult(tc.mockReturn)
			}

			// Call the actual UpdateCourse function
			actualReturn, err := s.service.UpdateCourse(context.Background(), tc.inputID, tc.inputCourse.Name)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedReturn, actualReturn)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestGetCourseByID() {
	t := s.T()

	course := models.Course{ID: 1, Name: "Databases"}

	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockRows       *sqlmock.Rows
		mockReturnErr  error
		inputID        int
		expectedReturn models.Course
		expectedError  error
	}{
		"course found by ID": {
			mockInputArgs:  []driver.Value{course.ID},
			mockRows:       sqlmock.NewRows([]string{"id", "name"}).AddRow(course.ID, course.Name),
			mockReturnErr:  nil,
			inputID:        course.ID,
			expectedReturn: course,
			expectedError:  nil,
		},
		"course not found": {
			mockInputArgs:  []driver.Value{999},
			mockRows:       sqlmock.NewRows([]string{"id", "name"}),
			mockReturnErr:  sql.ErrNoRows,
			inputID:        999,
			expectedReturn: models.Course{},
			expectedError:  fmt.Errorf("[in services.GetCourseByIDByID] no course found with id: %d", 999),
		},
		"Error retrieving course": {
			mockInputArgs:  []driver.Value{5},
			mockRows:       nil,
			mockReturnErr:  errors.New("test error"),
			inputID:        5,
			expectedReturn: models.Course{},
			expectedError:  fmt.Errorf("[in services.GetCourseByIDByID] failed to retrieve course: %w", errors.New("test error")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `SELECT id, name FROM course WHERE id = $1`
			mock := s.dbMock.ExpectQuery(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...)

			if tc.mockReturnErr != nil {
				mock.WillReturnError(tc.mockReturnErr)
			} else {
				mock.WillReturnRows(tc.mockRows)
			}

			actualReturn, err := s.service.GetCourseByID(context.Background(), tc.inputID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedReturn, actualReturn)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestCreateCourse() {
	t := s.T()

	courseName := "Databases"
	newID := 1

	testCases := map[string]struct {
		mockInputArgs  []driver.Value
		mockRows       *sqlmock.Rows
		mockReturnErr  error
		inputCourse    string
		expectedReturn models.Course
		expectedError  error
	}{
		"course created successfully": {
			mockInputArgs:  []driver.Value{courseName},
			mockRows:       sqlmock.NewRows([]string{"id"}).AddRow(newID), // Simulate returned new ID
			mockReturnErr:  nil,
			inputCourse:    courseName,
			expectedReturn: models.Course{ID: newID, Name: courseName},
			expectedError:  nil,
		},
		"error creating course": {
			mockInputArgs:  []driver.Value{courseName},
			mockRows:       nil,
			mockReturnErr:  errors.New("test error"),
			inputCourse:    courseName,
			expectedReturn: models.Course{},
			expectedError:  fmt.Errorf("failed to create course: %w", errors.New("test error")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			exp := `INSERT INTO course (name) VALUES ($1) RETURNING id`
			mock := s.dbMock.ExpectQuery(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...)

			if tc.mockReturnErr != nil {

				mock.WillReturnError(tc.mockReturnErr)
			} else {

				mock.WillReturnRows(tc.mockRows)
			}

			actualReturn, err := s.service.CreateCourse(context.Background(), tc.inputCourse)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedReturn, actualReturn)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func (s *testSuit) TestDeleteCourse() {
	t := s.T()

	testCases := map[string]struct {
		mockInputArgs []driver.Value
		mockReturn    driver.Result
		mockReturnErr error
		inputID       int
		expectedError error
	}{
		"course deleted successfully": {
			mockInputArgs: []driver.Value{1},
			mockReturn:    sqlmock.NewResult(0, 1),
			mockReturnErr: nil,
			inputID:       1,
			expectedError: nil,
		},
		"no course found with given ID": {
			mockInputArgs: []driver.Value{999},
			mockReturn:    sqlmock.NewResult(0, 0),
			mockReturnErr: nil,
			inputID:       999,
			expectedError: fmt.Errorf("[in services.ListCourses] no course found with id %d", 999),
		},
		"error executing delete": {
			mockInputArgs: []driver.Value{1},
			mockReturn:    nil,
			mockReturnErr: errors.New("test error"),
			inputID:       1,
			expectedError: fmt.Errorf("[in services.ListCourses] failed to delete course: %w", errors.New("test error")),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			exp := `DELETE FROM course WHERE id = $1`
			mock := s.dbMock.ExpectExec(regexp.QuoteMeta(exp)).
				WithArgs(tc.mockInputArgs...)

			if tc.mockReturnErr != nil {
				mock.WillReturnError(tc.mockReturnErr)
			} else {
				mock.WillReturnResult(tc.mockReturn)
			}

			err := s.service.DeleteCourse(context.Background(), tc.inputID)

			assert.Equal(t, tc.expectedError, err)

			err = s.dbMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
