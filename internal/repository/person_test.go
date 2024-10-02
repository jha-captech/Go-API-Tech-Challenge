package repository_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"

	// "jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/repository"
)

var testCourse1 = models.Course{
	ID:   11,
	Guid: "123",
	Name: "Test Course",
}

var testCourse2 = models.Course{
	ID:   12,
	Guid: "1234",
	Name: "Test Course2",
}

var testPerson = models.Person{
	ID:        100,
	Guid:      "100a",
	FirstName: "",
	LastName:  "",
	Email:     "test@test.com",
	Courses:   []models.Course{testCourse1, testCourse2},
}

func TestPersonFindAll(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "first_name", "last_name", "email", "age", "type"}).
		AddRow(100, "$$$", "Mr", "Krabs", "krustykrab@test.com", "48", "student").
		AddRow(101, "!!12", "Larry", "Lobster", "livinlikelarry@test.com", "31", "student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "person" WHERE first_name like $1 AND last_name like $2 AND email like $3`)).
		WithArgs("%M%", "%L%", "%com%").
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	people, _ := repo.FindAll(repository.Filters{
		"first_name": "M",
		"last_name":  "L",
		"email":      "com",
	})

	if len(people) != 2 {
		t.Errorf("expected 2 people returned. was: %d", len(people))
	}
}

func TestPersonFindAllNoFilters(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "first_name", "last_name", "email", "age", "type"}).
		AddRow(100, "$$$", "Mr", "Krabs", "krustykrab@test.com", "48", "student").
		AddRow(101, "!!12", "Larry", "Lobster", "livinlikelarry@test.com", "31", "student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "person"`)).
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	people, _ := repo.FindAll(repository.Filters{})

	if len(people) != 2 {
		t.Errorf("expected 2 people returned. was: %d", len(people))
	}
}

func TestPersonFindOne(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	rows3 := sqlmock.NewRows([]string{"Id", "guid", "first_name", "last_name", "email", "age", "type"}).
		AddRow(100, "100a", "Spongebob", "Squarepants", "bikinibottom@test.com", "55", "student")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "person" WHERE guid = $1`)).
		WithArgs("100a").
		WillReturnRows(rows3)

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "person_id", "course_id"}).
		AddRow(900, 100, 200).
		AddRow(901, 100, 201)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "person_course" WHERE "person_course"."person_id" = $1`)).
		WithArgs(100).
		WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"Id", "guid", "name"}).
		AddRow(200, "200", "Test 200").
		AddRow(201, "201", "Test 201")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course" WHERE "course"."id" IN ($1,$2)`)).
		WithArgs(200, 201).
		WillReturnRows(rows2)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	person, err := repo.FindOne("100a")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if person.Guid != "100a" ||
		person.FirstName != "Spongebob" ||
		person.LastName != "Squarepants" ||
		person.Age != 55 ||
		person.Email != "bikinibottom@test.com" ||
		person.Type != "student" {
		t.Errorf("Person record was not as expected. Was %v", person)
	}

}

func TestPersonFindOneNotFound(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	rows3 := sqlmock.NewRows([]string{"Id", "guid", "first_name", "last_name", "email", "age", "type"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "person" WHERE guid = $1`)).
		WithArgs("100a").
		WillReturnRows(rows3)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	_, err := repo.FindOne("100a")

	if fmt.Sprint(err) != fmt.Sprint(apperror.NotFound("Person: 100a Not Found")) {
		t.Errorf("Person Not Found Error not returned, was: %v", err)
	}
}

func TestDeletePerson(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE person_id = $1`)).
		WithArgs(testPerson.ID)
	mock.ExpectCommit()
	mock.ExpectClose()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM "person" WHERE "person"."id" = $1`)).
		WithArgs(testPerson.ID)
	mock.ExpectCommit()
	mock.ExpectClose()
	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	repo.Delete(&testPerson)
}
