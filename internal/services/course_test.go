package services

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
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
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			exp := `UPDATE course 
			SET name = $1 
			WHERE id = $2`
			log.Println(regexp.QuoteMeta(exp))
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
