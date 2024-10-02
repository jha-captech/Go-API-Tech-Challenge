package repository_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/repository"
)

func TestCourseFindAll(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "name"}).
		AddRow(100, "abcd", "Foo").
		AddRow(102, "1234", "Bar")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course" WHERE name like $1`)).
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	result, _ := repo.FindAll(repository.Filters{"name": "Foo"})

	if len(result) != 2 {
		t.Errorf("Expected 2 rows returned")
	}
}

func TestCourseFindAllNoFilters(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "name"}).
		AddRow(100, "abcd", "Foo").
		AddRow(102, "1234", "Bar")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course"`)).
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	result, _ := repo.FindAll(repository.Filters{})

	if len(result) != 2 {
		t.Errorf("Expected 2 rows returned")
	}
}

func TestCourseFindOne(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "name"}).
		AddRow(100, "abcd", "Foo")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course" WHERE guid = $1`)).
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	course, _ := repo.FindOne("123")

	if course.ID != 0 ||
		course.Name != "Foo" ||
		course.Guid != "abcd" {
		t.Errorf("Course Was not as expected was %v", course)
	}
}

func TestCourseFindOneNotFound(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "name"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "course" WHERE guid = $1`)).
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	_, err := repo.FindOne("123")

	if fmt.Sprint(err) != fmt.Sprint(apperror.NotFound("Course: 123 Not Found")) {
		t.Errorf("Not Found error not returned %v", err)
	}
}

func TestCourseSaveCreate(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	newCourse := models.Course{
		ID:   0,
		Guid: "abcd",
		Name: "Test Name",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "course" ("guid","name") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(newCourse.Guid, newCourse.Name).
		WillReturnRows(sqlmock.NewRows([]string{"Id", "guid", "name"}).
			AddRow(100, newCourse.Guid, newCourse.Name))
	mock.ExpectCommit()
	mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	err := repo.Save(&newCourse)

	if newCourse.Guid != "abcd" ||
		newCourse.Name != "Test Name" {
		t.Errorf("Course was not as expected was %v", newCourse)
	}
	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestCourseSaveUpdate(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	newCourse := models.Course{
		ID:   100,
		Guid: "abcd",
		Name: "Test Name",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "course" SET "guid"=$1,"name"=$2 WHERE "id" = $3`)).
		WithArgs(newCourse.Guid, newCourse.Name, newCourse.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	err := repo.Save(&newCourse)

	if newCourse.Guid != "abcd" ||
		newCourse.Name != "Test Name" {
		t.Errorf("Course was not as expected was %v", newCourse)
	}
	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestCourseDelete(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	testCourse := models.Course{
		ID:   100,
		Guid: "abcd",
		Name: "Test Name",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE course_id = $1`)).
		WithArgs(testCourse.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "course" WHERE "course"."id" = $1`)).
		WithArgs(testCourse.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewCourse(db, appLog)

	err := repo.Delete(&testCourse)

	if testCourse.Guid != "abcd" ||
		testCourse.Name != "Test Name" {
		t.Errorf("Course was not as expected was %v", testCourse)
	}

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}
