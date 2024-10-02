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
		// Unfortunately filters can be applied in a random order so when including all filters, the order is non deterministic.
	mock.ExpectQuery(`SELECT (.+)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	people, err := repo.FindAll(repository.Filters{
		"first_name": "M",
		"last_name":  "L",
		"email":      "com",
	})

	if len(people) != 2 {
		t.Errorf("expected 2 people returned. was: %d", len(people))
	}

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestPersonFindAllOneFilter(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	// Define Mock Database behavior
	rows := sqlmock.NewRows([]string{"Id", "guid", "first_name", "last_name", "email", "age", "type"}).
		AddRow(100, "$$$", "Mr", "Krabs", "krustykrab@test.com", "48", "student").
		AddRow(101, "!!12", "Larry", "Lobster", "livinlikelarry@test.com", "31", "student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "person" WHERE first_name like $1`)).
		WithArgs("%M%").
		WillReturnRows(rows)

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	people, err := repo.FindAll(repository.Filters{
		"first_name": "M",
	})

	if len(people) != 2 {
		t.Errorf("expected 2 people returned. was: %d", len(people))
	}

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
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

	people, err := repo.FindAll(repository.Filters{})

	if len(people) != 2 {
		t.Errorf("expected 2 people returned. was: %d", len(people))
	}

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
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
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE person_id = $1`)).
		WithArgs(testPerson.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person" WHERE "person"."id" = $1`)).
		WithArgs(testPerson.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectClose()
	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	err := repo.Delete(&testPerson)
	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestSaveNewPerson(t *testing.T) {
	newPerson := models.Person{
		// ID:        0,
		Guid:      "0000",
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@test.com",
		Age:       30,
	}

	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "person" ("guid","first_name","last_name","email","age","type") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs("0000", newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.Age, newPerson.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "course" ("guid","name","id") VALUES ($1,$2,$3),($4,$5,$6) ON CONFLICT DO NOTHING RETURNING "id"`)).
		WithArgs(testCourse1.Guid, testCourse1.Name, testCourse1.ID, testCourse2.Guid, testCourse2.Name, testCourse2.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testCourse1.ID).AddRow(testCourse2.ID))
		// WithArgs(1, testCourse1, 1, testCourse2)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "person_course" ("person_id","course_id") VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING`)).
		WithArgs(1, testCourse1.ID, 1, testCourse2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	err := repo.Save(&newPerson, []models.Course{testCourse1, testCourse2})

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestSaveUpdatedPerson(t *testing.T) {
	newPerson := models.Person{
		ID:        1,
		Guid:      "0000",
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@test.com",
		Age:       30,
	}

	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "person" SET "guid"=$1,"first_name"=$2,"last_name"=$3,"email"=$4,"age"=$5,"type"=$6 WHERE "id" = $7`)).
		WithArgs("0000", newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.Age, newPerson.Type, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
		// WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	mock.ExpectClose()

	// Add Person_course

	// mock.ExpectBegin()

	// mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "person_course" ("person_id","course_id") VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING`)).
	// 	WithArgs(1, testCourse1, 1, testCourse2)

	// mock.ExpectCommit()
	// mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	err := repo.Save(&newPerson, nil)

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestSaveUpdateRemoveCourse(t *testing.T) {
	newPerson := models.Person{
		ID:        1,
		Guid:      "0000",
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@test.com",
		Age:       30,
		Courses:   []models.Course{testCourse1, testCourse2},
	}

	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "person_course" WHERE ("person_course"."person_id","person_course"."course_id") IN (($1,$2)`)).
		WithArgs(1, testCourse2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "person" SET "guid"=$1,"first_name"=$2,"last_name"=$3,"email"=$4,"age"=$5,"type"=$6 WHERE "id" = $7`)).
		WithArgs("0000", newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.Age, newPerson.Type, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "course" ("guid","name","id") VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "id"`)).
		WithArgs(testCourse1.Guid, testCourse1.Name, testCourse1.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(testCourse1.ID))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "person_course" ("person_id","course_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`)).
		WithArgs(1, testCourse1.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	err := repo.Save(&newPerson, []models.Course{testCourse1})

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}

func TestSaveUpdateAddCourse(t *testing.T) {
	newPerson := models.Person{
		ID:        1,
		Guid:      "0000",
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@test.com",
		Age:       30,
		Courses:   []models.Course{testCourse1},
	}

	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "person" SET "guid"=$1,"first_name"=$2,"last_name"=$3,"email"=$4,"age"=$5,"type"=$6 WHERE "id" = $7`)).
		WithArgs("0000", newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.Age, newPerson.Type, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "course" ("guid","name","id") VALUES ($1,$2,$3),($4,$5,$6) ON CONFLICT DO NOTHING RETURNING "id"`)).
		WithArgs(testCourse1.Guid, testCourse1.Name, testCourse1.ID, testCourse2.Guid, testCourse2.Name, testCourse2.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(testCourse1.ID))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "person_course" ("person_id","course_id") VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING`)).
		WithArgs(1, testCourse1.ID, 1, testCourse2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectClose()

	db, _ := gorm.Open(dialector, &gorm.Config{})

	appLog := applog.New(log.Default())

	repo := repository.NewPerson(db, appLog)

	err := repo.Save(&newPerson, []models.Course{testCourse1, testCourse2})

	if err != nil {
		t.Errorf("Error was returned when not expected %v", err)
	}
}
