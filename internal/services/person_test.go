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

var testCourse1 = models.Course{
	ID:   1,
	Guid: "1",
	Name: "Test 1",
}

var testCourse2 = models.Course{
	ID:   2,
	Guid: "2",
	Name: "Test 2",
}

var testPerson = models.Person{
	ID:        11,
	Guid:      "1234",
	FirstName: "Spongebob",
	LastName:  "Squarepants",
	Email:     "spants@test.com",
	Age:       45,
	Type:      string(models.Student),
	Courses:   []models.Course{testCourse1},
}

var personGetOneByGuidTestCase = []struct {
	name           string
	guid           string
	expectedPerson models.Person
	expectedErr    error
}{
	{name: "Success Case", guid: "1234", expectedPerson: testPerson, expectedErr: nil},
	{name: "Person Not Found", guid: "abcd", expectedPerson: models.Person{}, expectedErr: apperror.NotFound("Person: abcd Not Found")},
}

func TestPersonGetOneByGuid(t *testing.T) {
	mockPersonRepo := new(jfmock.Person)
	mockCourseRepo := new(jfmock.Course)

	mockPersonRepo.On("FindOne", "1234").Return(testPerson, nil)
	mockPersonRepo.On("FindOne", "abcd").Return(models.Person{}, apperror.NotFound("Person: abcd Not Found"))

	personService := services.NewPerson(&applog.AppLogger{}, mockPersonRepo, mockCourseRepo)

	for _, tc := range personGetOneByGuidTestCase {
		t.Run(tc.name, func(t *testing.T) {
			outPerson, outErr := personService.GetOneByGuid(tc.guid)

			if fmt.Sprint(outPerson) != fmt.Sprint(tc.expectedPerson) {
				t.Errorf("Returned person was not as expected. want: %v got: %v", tc.expectedPerson, outPerson)
			}

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as epxected. want: %v got: %v", tc.expectedErr, outErr)
			}
		})
	}
}

var personUpdateTestCase = []struct {
	name            string
	guid            string
	input           models.PersonInput
	expectedPerson  models.Person
	expectedCourses []models.Course
	expectedErr     error
	onRepoUpdate    error
}{
	{
		name: "Success Case",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: []string{"1", "2"},
		},
		expectedPerson: models.Person{
			ID:        testPerson.ID,
			Guid:      testPerson.Guid,
			FirstName: "Patrick",
			LastName:  "Star",
			Email:     "pstar@test.com",
			Age:       12,
			Type:      string(models.Student),
			Courses:   []models.Course{testCourse1},
		},
		expectedCourses: []models.Course{testCourse1, testCourse2},
	},
	{
		name: "Success Case At Least 10 years old",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         10,
			Type:        string(models.Student),
			CourseGuids: []string{"1", "2"},
		},
		expectedPerson: models.Person{
			ID:        testPerson.ID,
			Guid:      testPerson.Guid,
			FirstName: "Patrick",
			LastName:  "Star",
			Email:     "pstar@test.com",
			Age:       10,
			Type:      string(models.Student),
			Courses:   []models.Course{testCourse1},
		},
		expectedCourses: []models.Course{testCourse1, testCourse2},
	},
	{
		name: "Success Update Course",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: []string{"2"},
		},
		expectedPerson: models.Person{
			ID:        testPerson.ID,
			Guid:      testPerson.Guid,
			FirstName: "Patrick",
			LastName:  "Star",
			Email:     "pstar@test.com",
			Age:       12,
			Type:      string(models.Student),
			Courses:   []models.Course{testCourse1},
		},
		expectedCourses: []models.Course{testCourse2},
	},

	{
		name: "Success Remove All Courses",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: models.Person{
			ID:        testPerson.ID,
			Guid:      testPerson.Guid,
			FirstName: "Patrick",
			LastName:  "Star",
			Email:     "pstar@test.com",
			Age:       12,
			Type:      string(models.Student),
			Courses:   []models.Course{testCourse1},
		},
		expectedCourses: nil,
	},
	{
		name: "Not Found",
		guid: "abcd",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: []string{"1", "2"},
		},
		expectedPerson:  models.Person{},
		expectedCourses: []models.Course{testCourse1, testCourse2},
		expectedErr:     apperror.BadRequest("Person: abcd Not Found"),
	},
	{
		name: "Repository Error",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: models.Person{
			ID:        testPerson.ID,
			Guid:      testPerson.Guid,
			FirstName: "Patrick",
			LastName:  "Star",
			Email:     "pstar@test.com",
			Age:       12,
			Type:      string(models.Student),
			Courses:   []models.Course{testCourse1},
		},
		expectedErr:  apperror.BadRequest("Repo Error"),
		onRepoUpdate: apperror.BadRequest("Repo Error"),
	},
	{
		name: "First Name Not Blank",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   " ",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr:    apperror.BadRequest("First Name must not be blank"),
	},
	{
		name: "Last Name Not Blank",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patric",
			LastName:    "",
			Email:       "pstar@test.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr:    apperror.BadRequest("Last Name must not be blank"),
	},
	{
		name: "Email Not Valid",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "dwdetest.com",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr:    apperror.BadRequest("Email must be a valid email address"),
	},
	{
		name: "Email Not Valid Blank",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "",
			Age:         12,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr:    apperror.BadRequest("Email must be a valid email address"),
	},

	{
		name: "Age must be at least 10",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         9,
			Type:        string(models.Student),
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr:    apperror.BadRequest("Must be at least 10 years old to enrol."),
	},

	{
		name: "Type invalid",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   "Patrick",
			LastName:    "Star",
			Email:       "pstar@test.com",
			Age:         120,
			Type:        "Foo",
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr:    apperror.BadRequest("Invalid Person type, must be either 'professor' or 'student'"),
	},
	{
		name: "Multi Error",
		guid: "1234",
		input: models.PersonInput{
			FirstName:   " ",
			LastName:    " ",
			Email:       " ",
			Age:         1,
			Type:        "Foo",
			CourseGuids: nil,
		},
		expectedPerson: testPerson,
		expectedErr: apperror.Of([]error{
			apperror.BadRequest("First Name must not be blank"),
			apperror.BadRequest("Last Name must not be blank"),
			apperror.BadRequest("Email must be a valid email address"),
			apperror.BadRequest("Must be at least 10 years old to enrol."),
			apperror.BadRequest("Invalid Person type, must be either 'professor' or 'student'"),
		}),
	},
}

func TestPersonUpdate(t *testing.T) {

	for _, tc := range personUpdateTestCase {
		t.Run(tc.name, func(t *testing.T) {

			mockPersonRepo := new(jfmock.Person)
			mockCourseRepo := new(jfmock.Course)

			mockPersonRepo.On("FindOne", "1234").Return(testPerson, nil)
			mockPersonRepo.On("FindOne", "abcd").Return(models.Person{}, apperror.NotFound("Person: abcd Not Found"))

			mockCourseRepo.On("FindOne", "1").Return(testCourse1, nil)
			mockCourseRepo.On("FindOne", "2").Return(testCourse2, nil)
			mockCourseRepo.On("FindOne", "3").Return(models.Course{}, apperror.NotFound("Course Not Found"))

			mockPersonRepo.On("Save", &tc.expectedPerson, tc.expectedCourses).Return(tc.onRepoUpdate)

			personService := services.NewPerson(&applog.AppLogger{}, mockPersonRepo, mockCourseRepo)

			outPerson, outErr := personService.Update(tc.guid, tc.input)

			if fmt.Sprint(outPerson) != fmt.Sprint(tc.expectedPerson) {
				t.Errorf("Returned person was not as expected. want: %v got: %v", tc.expectedPerson, outPerson)
			}

			// outMultiErr, isOutMulti := outErr.(apperror.MultiError)
			// expectedMultiErr, isExpectMulti := tc.expectedErr.(apperror.MultiError)

			// if isOutMulti || isExpectMulti {
			// 	t.Errorf("Out: %v Multi: %v \n", outMultiErr.Errors, expectedMultiErr)
			// }

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as epxected. want: %v got: %v", tc.expectedErr, outErr)
			}
		})
	}
}

var personDeleteTestCase = []struct {
	name           string
	guid           string
	expectedPerson models.Person
	expectedErr    error
	repoDeleteErr  error
}{
	{
		name:           "Success",
		guid:           "1234",
		expectedPerson: testPerson,
		expectedErr:    nil,
	},
	{
		name:           "Person Not Found",
		guid:           "abcd",
		expectedPerson: models.Person{},
		expectedErr:    apperror.NotFound("Person: abcd Not Found"),
	},
	{
		name:           "Repo Error",
		guid:           "1234",
		expectedPerson: testPerson,
		expectedErr:    apperror.NotFound("Repo Error"),
		repoDeleteErr:  apperror.BadRequest("Repo Error"),
	},
}

func TestPersonDelete(t *testing.T) {

	for _, tc := range personDeleteTestCase {
		t.Run(tc.name, func(t *testing.T) {
			mockPersonRepo := new(jfmock.Person)
			mockCourseRepo := new(jfmock.Course)

			mockPersonRepo.On("FindOne", "1234").Return(testPerson, nil)
			mockPersonRepo.On("FindOne", "abcd").Return(models.Person{}, apperror.NotFound("Person: abcd Not Found"))

			mockPersonRepo.On("Delete", &tc.expectedPerson).Return(tc.repoDeleteErr)
			personService := services.NewPerson(&applog.AppLogger{}, mockPersonRepo, mockCourseRepo)

			outErr := personService.Delete(tc.guid)

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as expected. want: %v got: %v", tc.expectedErr, outErr)
			}
		})
	}
}

var personCreateTestCase = []struct {
	name            string
	input           models.PersonInput
	expectedPerson  models.Person
	expectedCourses []models.Course
	expectedErr     error
	onRepoSave      error
}{
	{
		name: "Success Case",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         10,
			Type:        string(models.Professor),
			CourseGuids: []string{"1", "2"},
		},
		expectedPerson: models.Person{
			ID:        1,
			Guid:      "",
			FirstName: "Squidward",
			LastName:  "Tentacles",
			Email:     "boldandbrash@test.com",
			Age:       10,
			Type:      string(models.Professor),
			Courses:   []models.Course{testCourse1, testCourse2},
		},
		expectedCourses: []models.Course{testCourse1, testCourse2},
	},
	{
		name: "Success Case No Courses",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         55,
			Type:        string(models.Professor),
			CourseGuids: nil,
		},
		expectedPerson: models.Person{
			ID:        1,
			Guid:      "",
			FirstName: "Squidward",
			LastName:  "Tentacles",
			Email:     "boldandbrash@test.com",
			Age:       55,
			Type:      string(models.Professor),
			Courses:   nil,
		},
		expectedCourses: nil,
	},
	{
		name: "Repo Error",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         55,
			Type:        string(models.Professor),
			CourseGuids: []string{"1", "2"},
		},
		expectedPerson: models.Person{
			ID:        1,
			Guid:      "",
			FirstName: "Squidward",
			LastName:  "Tentacles",
			Email:     "boldandbrash@test.com",
			Age:       55,
			Type:      string(models.Professor),
			Courses:   []models.Course{testCourse1, testCourse2},
		},
		expectedCourses: []models.Course{testCourse1, testCourse2},
		onRepoSave:      apperror.BadRequest("Repo Error"),
		expectedErr:     apperror.BadRequest("Repo Error"),
	},
	{
		name: "Course Not Found",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         55,
			Type:        string(models.Professor),
			CourseGuids: []string{"3"},
		},
		expectedPerson: models.Person{
			ID:        1,
			Guid:      "",
			FirstName: "Squidward",
			LastName:  "Tentacles",
			Email:     "boldandbrash@test.com",
			Age:       55,
			Type:      string(models.Professor),
		},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("Course Not Found"),
	},
	{
		name: "First Name Blank",
		input: models.PersonInput{
			FirstName:   " ",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         55,
			Type:        string(models.Professor),
			CourseGuids: []string{},
		},
		expectedPerson:  models.Person{},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("First Name must not be blank"),
	},
	{
		name: "Last Name Blank",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    " ",
			Email:       "boldandbrash@test.com",
			Age:         55,
			Type:        string(models.Professor),
			CourseGuids: []string{},
		},
		expectedPerson:  models.Person{},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("Last Name must not be blank"),
	},
	{
		name: "Email Valid",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrashtest.com",
			Age:         55,
			Type:        string(models.Professor),
			CourseGuids: []string{},
		},
		expectedPerson:  models.Person{},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("Email must be a valid email address"),
	},
	{
		name: "Age less than 10",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         9,
			Type:        string(models.Professor),
			CourseGuids: []string{},
		},
		expectedPerson:  models.Person{},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("Must be at least 10 years old to enrol."),
	},
	{
		name: "Invalid Person Type",
		input: models.PersonInput{
			FirstName:   "Squidward",
			LastName:    "Tentacles",
			Email:       "boldandbrash@test.com",
			Age:         90,
			Type:        "Foo",
			CourseGuids: []string{},
		},
		expectedPerson:  models.Person{},
		expectedCourses: nil,
		expectedErr:     apperror.BadRequest("Invalid Person type, must be either 'professor' or 'student'"),
	},
	{
		name: "Multi Error",
		input: models.PersonInput{
			FirstName:   " ",
			LastName:    " ",
			Email:       "boldandbrash",
			Age:         9,
			Type:        "Foo",
			CourseGuids: []string{},
		},
		expectedPerson:  models.Person{},
		expectedCourses: nil,
		expectedErr: apperror.Of([]error{
			apperror.BadRequest("First Name must not be blank"),
			apperror.BadRequest("Last Name must not be blank"),
			apperror.BadRequest("Email must be a valid email address"),
			apperror.BadRequest("Must be at least 10 years old to enrol."),
			apperror.BadRequest("Invalid Person type, must be either 'professor' or 'student'"),
		}),
	},
}

func TestPersonCreate(t *testing.T) {
	for _, tc := range personCreateTestCase {
		mockPersonRepo := new(jfmock.Person)
		mockCourseRepo := new(jfmock.Course)

		mockPersonRepo.On("FindOne", "1234").Return(testPerson, nil)
		mockPersonRepo.On("FindOne", "abcd").Return(models.Person{}, apperror.NotFound("Person: abcd Not Found"))

		mockCourseRepo.On("FindOne", "1").Return(testCourse1, nil)
		mockCourseRepo.On("FindOne", "2").Return(testCourse2, nil)
		mockCourseRepo.On("FindOne", "3").Return(models.Course{}, apperror.NotFound("Course Not Found"))

		mockPersonRepo.On("Save", mock.Anything, tc.expectedCourses).Return(tc.onRepoSave)

		personService := services.NewPerson(&applog.AppLogger{}, mockPersonRepo, mockCourseRepo)

		outPerson, outErr := personService.Create(tc.input)

		if outPerson.FirstName != tc.expectedPerson.FirstName ||
			outPerson.LastName != tc.expectedPerson.LastName ||
			outPerson.Email != tc.expectedPerson.Email ||
			outPerson.Age != tc.expectedPerson.Age ||
			outPerson.Type != tc.expectedPerson.Type {
			t.Errorf("Returned person was not as expected. want: %v got: %v", tc.expectedPerson, outPerson)
		}

		if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
			t.Errorf("Returned error was not as epxected. want: %v got: %v", tc.expectedErr, outErr)
		}
	}
}

var personGetAllTestCase = []struct {
	name            string
	urlParam        url.Values
	expectedPersons []models.Person
	expectedFilters repository.Filters
	expectedErr     error
	repoErr         error
}{
	{
		name:            "Success",
		expectedPersons: []models.Person{testPerson, testPerson},
		expectedFilters: repository.Filters{},
	},
	{
		name:            "Repo Error",
		expectedPersons: []models.Person{},
		expectedFilters: repository.Filters{},
		expectedErr:     apperror.BadRequest("Repo Error"),
		repoErr:         apperror.BadRequest("Repo Error"),
	},
	{
		name:            "Success All Filters",
		expectedPersons: []models.Person{testPerson, testPerson},
		urlParam: url.Values{
			"FirstName": []string{"Foo"},
			"LastName":  []string{"Bar"},
			"Email":     []string{"foobar"},
		},
		expectedFilters: repository.Filters{
			"first_name": "Foo",
			"last_name":  "Bar",
			"email":      "foobar",
		},
	},
	{
		name:            "Fail Invalid filter",
		expectedPersons: nil,
		urlParam: url.Values{
			"FirstNa":  []string{"Foo"},
			"LastName": []string{"Bar"},
			"Email":    []string{"foobar"},
		},
		expectedFilters: repository.Filters{},
		expectedErr:     apperror.BadRequest("Invalid Request Parameter: FirstNa"),
	},
	{
		name:            "Fail Multi Error",
		expectedPersons: nil,
		urlParam: url.Values{
			"FirstNa": []string{"Foo"},
			"LastNa":  []string{"Bar"},
			"Email":   []string{"foobar"},
		},
		expectedFilters: repository.Filters{},
		expectedErr: apperror.Of([]error{
			apperror.BadRequest("Invalid Request Parameter: FirstNa"),
			apperror.BadRequest("Invalid Request Parameter: LastNa"),
		}),
	},
}

func TestPersonGetAll(t *testing.T) {
	for _, tc := range personGetAllTestCase {
		t.Run(tc.name, func(t *testing.T) {
			mockPersonRepo := new(jfmock.Person)
			mockCourseRepo := new(jfmock.Course)

			mockPersonRepo.On("FindAll", tc.expectedFilters).Return(tc.expectedPersons, tc.repoErr)

			personService := services.NewPerson(&applog.AppLogger{}, mockPersonRepo, mockCourseRepo)

			outPersons, outErr := personService.GetAll(tc.urlParam)

			if len(outPersons) != len(tc.expectedPersons) {
				t.Errorf("Length was not matched, want %d got %d", len(tc.expectedPersons), len(outPersons))
			}

			if fmt.Sprint(outErr) != fmt.Sprint(tc.expectedErr) {
				t.Errorf("Returned error was not as epxected. want: %v got: %v", tc.expectedErr, outErr)
			}
		})
	}
}
