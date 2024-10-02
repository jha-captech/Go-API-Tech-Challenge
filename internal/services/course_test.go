package services_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/repository"
	jfmock "jf.go.techchallenge/internal/repository/mock"
	"jf.go.techchallenge/internal/services"
)

var testCourse = models.Course{
	ID:   11,
	Guid: "123",
	Name: "Test Course",
}

var getOneByGuidTestCase = []struct {
	name           string
	guid           string
	expectedCourse models.Course
	expectedErr    error
}{
	{name: "Sucess Case", guid: "123", expectedCourse: testCourse, expectedErr: nil},
	{name: "Course Not Found", guid: "abcd", expectedCourse: models.Course{}, expectedErr: apperror.NotFound("Course: abcd Not Found")},
}

func Test_Course_GetOneByGuid(t *testing.T) {
	mockCourseRepo := new(jfmock.Course)

	mockCourseRepo.On("FindOne", "123").Return(testCourse, nil)
	mockCourseRepo.On("FindOne", "abcd").Return(models.Course{}, apperror.NotFound("Course: abcd Not Found"))

	courseService := services.NewCourse(&applog.AppLogger{}, mockCourseRepo)

	for _, tc := range getOneByGuidTestCase {
		t.Run(tc.name, func(t *testing.T) {

			outCourse, outErr := courseService.GetOneByGuid(tc.guid)

			if outCourse != tc.expectedCourse {
				t.Errorf("Returned course was not as expected. want: %v got: %v", tc.expectedCourse, outCourse)
			}

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as epxected. want: %v got: %v", tc.expectedErr, outErr)
			}
		})
	}
}

var updateTestCase = []struct {
	name           string
	guid           string
	input          models.CourseInput
	expectedCourse models.Course
	expectedErr    error
	onRepoUpdate   error
}{
	{
		name:  "Success",
		guid:  "123",
		input: models.CourseInput{Name: "Test Update"},
		expectedCourse: models.Course{
			ID:   testCourse.ID,
			Name: "Test Update",
			Guid: testCourse.Guid,
		},
		expectedErr:  nil,
		onRepoUpdate: nil,
	},
	{
		name:           "Not Found",
		guid:           "abcd",
		input:          models.CourseInput{Name: "Test Update"},
		expectedCourse: models.Course{},
		expectedErr:    apperror.NotFound("Course: abcd Not Found"),
		onRepoUpdate:   nil,
	},
	{
		name:           "Name Blank",
		guid:           "123",
		input:          models.CourseInput{Name: " "},
		expectedCourse: testCourse, // TODO is this the best behavior? In an error case should the blank object models.CourseInput be returned?
		// https://google.github.io/styleguide/go/decisions#errors
		expectedErr:  apperror.BadRequest("Name must not be blank"),
		onRepoUpdate: nil,
	},
	{
		name:  "Repository Error",
		guid:  "123",
		input: models.CourseInput{Name: "Repo Fail"},
		expectedCourse: models.Course{
			ID:   testCourse.ID,
			Name: "Repo Fail",
			Guid: testCourse.Guid,
		},
		expectedErr:  apperror.BadRequest("Repository Fail"),
		onRepoUpdate: apperror.BadRequest("Repository Fail"),
	},
}

func Test_Update(t *testing.T) {
	mockCourseRepo := new(jfmock.Course)

	mockCourseRepo.On("FindOne", "123").Return(testCourse, nil)
	mockCourseRepo.On("FindOne", "abcd").Return(models.Course{}, apperror.NotFound("Course: abcd Not Found"))

	courseService := services.NewCourse(&applog.AppLogger{}, mockCourseRepo)

	for _, tc := range updateTestCase {
		t.Run(tc.name, func(t *testing.T) {

			mockCourseRepo.On("Save", &tc.expectedCourse).Return(tc.onRepoUpdate)
			outCourse, outErr := courseService.Update(tc.guid, tc.input)

			if outCourse != tc.expectedCourse {
				t.Errorf("Returned course was not as expected. want: %v got: %v", tc.expectedCourse, outCourse)
			}

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as expected. want: %v got: %v", tc.expectedErr, outErr)
			}
		})
	}
}

var deleteTestCase = []struct {
	name           string
	guid           string
	expectedCourse models.Course
	expectedErr    error
	onRepoDelete   error
}{
	{name: "Success Case", guid: "123", expectedCourse: testCourse, expectedErr: nil, onRepoDelete: nil},
	{name: "Course Not Found", guid: "abcd", expectedCourse: models.Course{}, expectedErr: apperror.NotFound("Course: abcd Not Found")},
	{name: "Repo Delete Error", guid: "123", expectedCourse: testCourse, expectedErr: apperror.BadRequest("Repo Error"), onRepoDelete: apperror.BadRequest("Repo Error")},
}

func Test_Delete(t *testing.T) {

	for _, tc := range deleteTestCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCourseRepo := new(jfmock.Course)

			mockCourseRepo.On("FindOne", "123").Return(testCourse, nil)
			mockCourseRepo.On("FindOne", "abcd").Return(models.Course{}, apperror.NotFound("Course: abcd Not Found"))

			courseService := services.NewCourse(&applog.AppLogger{}, mockCourseRepo)

			mockCourseRepo.On("Delete", &tc.expectedCourse).Return(tc.onRepoDelete)

			outErr := courseService.Delete(tc.guid)

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as expected. want: %v got: %v, return: %v", tc.expectedErr, outErr, tc.onRepoDelete)
			}

		})
	}
}

var createTestCase = []struct {
	name           string
	input          models.CourseInput
	expectedErr    error
	expectedCourse models.Course
	onRepoCreate   error
}{
	{
		name: "Success",
		input: models.CourseInput{
			Name: "Test Course",
		},
		expectedErr:    nil,
		expectedCourse: models.Course{Guid: "Any", ID: 0, Name: "Test Course"},
		onRepoCreate:   nil,
	},
	{
		name: "Name Blank",
		input: models.CourseInput{
			Name: "     ",
		},
		expectedErr:  apperror.BadRequest("Name must not be blank"),
		onRepoCreate: nil,
	},
	{
		name:           "Repo Save Fail",
		onRepoCreate:   apperror.BadRequest("Repo Fail"),
		expectedErr:    apperror.BadRequest("Repo Fail"),
		input:          models.CourseInput{Name: "Test Course"},
		expectedCourse: models.Course{Name: "Test Course"},
	},
}

func Test_Create(t *testing.T) {
	for _, tc := range createTestCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCourseRepo := new(jfmock.Course)

			courseService := services.NewCourse(&applog.AppLogger{}, mockCourseRepo)

			mockCourseRepo.On("Save", mock.Anything).Return(tc.onRepoCreate)

			outCourse, outErr := courseService.Create(tc.input)

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as expected. want: %v got: %v, return: %v", tc.expectedErr, outErr, tc.onRepoCreate)
			}

			// Courses will have guids generated on create so they will need to be asserted manually.
			if outCourse.ID != tc.expectedCourse.ID ||
				outCourse.Name != tc.expectedCourse.Name {
				t.Errorf("Returned course was not as expected. want: %v got: %v", tc.expectedCourse, outCourse)
			}

		})
	}
}

var getAllTestCase = []struct {
	name            string
	urlParams       url.Values
	expectedCourses []models.Course
	expectedFilters repository.Filters
	expectedErr     error
	onRepoErr       error
}{
	{
		name:            "Success",
		urlParams:       url.Values{"Name": []string{"foo"}},
		expectedCourses: []models.Course{testCourse, testCourse},
		expectedFilters: repository.Filters{"name": "foo"},
		expectedErr:     nil,
		onRepoErr:       nil,
	},
	{
		name:            "Success No filters",
		urlParams:       nil,
		expectedCourses: []models.Course{testCourse, testCourse, testCourse},
		expectedFilters: repository.Filters{},
		expectedErr:     nil,
		onRepoErr:       nil,
	},
	{
		name:            "Invalid Filter",
		urlParams:       url.Values{"Foo": []string{"foo"}},
		expectedFilters: repository.Filters{},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("Invalid Request Parameter: Foo"),
		onRepoErr:       nil,
	},
}

func Test_GetAll(t *testing.T) {
	for _, tc := range getAllTestCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCourseRepo := new(jfmock.Course)
			mockCourseRepo.On("FindAll", tc.expectedFilters).Return(tc.expectedCourses, tc.onRepoErr)

			courseService := services.NewCourse(&applog.AppLogger{}, mockCourseRepo)

			courses, outErr := courseService.GetAll(tc.urlParams)

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as expected. want: %v got: %v, return: %v", tc.expectedErr, outErr, tc.onRepoErr)
			}

			if len(courses) != len(tc.expectedCourses) {
				t.Errorf("Expected number of courses were not returned wanted %d got %d", len(tc.expectedCourses), len(courses))
			}
		})
	}
}
