package services

import (
	"context"
	"database/sql"
	"fmt"
	"go-api-tech-challenge/internal/models"

	"github.com/lib/pq"
)

type PersonService struct {
	database *sql.DB
}

// NewUserService returns a new UserService struct.
func NewPersonService(db *sql.DB) *PersonService {
	return &PersonService{
		database: db,
	}
}

func (s *PersonService) ListPersons(ctx context.Context) ([]models.Person, error) {

	query := `SELECT p.id as person_id, 
	p.first_name, 
	p.last_name,
	p.type,
	p.age,
	COALESCE(Array_AGG(pc.course_id), '{}') as course_ids
	FROM person p
	LEFT JOIN person_course pc ON p.id = pc.person_id
	GROUP BY p.id, p.first_name, p.last_name, p.type, p.age
	ORDER BY person_id asc`
	rows, err := s.database.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return []models.Person{}, fmt.Errorf("[in services.ListPersons] failed to get persons: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		var dbCourseIDs []sql.NullInt64
		err = rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age, pq.Array(&dbCourseIDs))
		if err != nil {
			return []models.Person{}, fmt.Errorf("[in services.ListPersons] failed to scan person from row: %w", err)
		}
		// Convert pg array into []int
		intCourseIDs := make([]int, len(dbCourseIDs))
		if len(dbCourseIDs) > 0 && dbCourseIDs[0].Valid {
			for i, id := range dbCourseIDs {
				intCourseIDs[i] = int(id.Int64) // Convert int64 to int
			}
			person.Courses = intCourseIDs
		}
		persons = append(persons, person)
	}

	if err = rows.Err(); err != nil {
		return []models.Person{}, fmt.Errorf("[in services.ListCourses] failed to scan courses: %w", err)
	}

	return persons, nil

}

func (s *PersonService) GetPersonByName(ctx context.Context, name string) (models.Person, error) {
	var person models.Person
	var dbCourseIDs []sql.NullInt64
	query := `SELECT p.id as person_id, 
	p.first_name, 
	p.last_name,
	p.type,
	p.age,
	COALESCE(Array_AGG(pc.course_id), '{}') as course_ids
	FROM person p
	LEFT JOIN person_course pc ON p.id = pc.person_id
	WHERE LOWER(p.last_name) = LOWER($1)
	GROUP BY p.id, p.first_name, p.last_name, p.type, p.age`

	err := s.database.QueryRowContext(ctx, query, name).Scan(&person.ID, &person.FirstName,
		&person.LastName, &person.Type, &person.Age, pq.Array(&dbCourseIDs))
	intCourseIDs := make([]int, len(dbCourseIDs))
	if len(dbCourseIDs) > 0 && dbCourseIDs[0].Valid {
		for i, id := range dbCourseIDs {
			intCourseIDs[i] = int(id.Int64) // Convert int64 to int
		}
		person.Courses = intCourseIDs
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, fmt.Errorf("[in services.GetPersonByName] no person found with name: %s", name)
		}
		return models.Person{}, fmt.Errorf("[in services.GetPersonByName] failed to retrieve person: %w", err)
	}
	return person, nil
}

func (s *PersonService) UpdatePerson(ctx context.Context, lastName string, updatedPerson models.Person) (models.Person, error) {
	var person models.Person
	var updatedCourseIDs []int

	tx, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `
	UPDATE person
	SET first_name = $1,
	last_name = $2,
	type = $3,
	age = $4
	WHERE LOWER(last_name) = LOWER($5)
	RETURNING id, first_name, last_name, type, age;
	`

	err = s.database.QueryRow(query, updatedPerson.FirstName, updatedPerson.LastName,
		updatedPerson.Type, updatedPerson.Age, lastName).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.Type,
		&person.Age,
	)

	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("failed to update person: %w", err)
	}

	if len(updatedPerson.Courses) > 0 {

		deleteCoursesQuery := `DELETE FROM person_course WHERE person_id = $1`
		_, err = tx.Exec(deleteCoursesQuery, person.ID)
		if err != nil {
			tx.Rollback()
			return models.Person{}, fmt.Errorf("failed to delete existing courses: %w", err)
		}
		insertCoursesQuery := `INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)`
		for _, courseID := range updatedPerson.Courses {
			_, err := tx.Exec(insertCoursesQuery, person.ID, courseID)
			if err != nil {
				tx.Rollback()
				return models.Person{}, fmt.Errorf("failed to insert new courses: %w", err)
			}
		}

	}

	// Get course IDs from person_course
	selectCoursesQuery := `SELECT course_id FROM person_course WHERE person_id = $1`
	rows, err := tx.Query(selectCoursesQuery, person.ID)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("failed to retrieve updated courses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var courseID int
		if err := rows.Scan(&courseID); err != nil {
			tx.Rollback()
			return models.Person{}, fmt.Errorf("failed to scan course ID: %w", err)
		}
		updatedCourseIDs = append(updatedCourseIDs, courseID)
	}
	person.Courses = updatedCourseIDs
	err = tx.Commit()
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return person, nil

}
func (s *PersonService) CreatePerson(ctx context.Context, person models.Person) (models.Person, error) {
	tx, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO person (first_name, last_name, type, age) 
	          VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, type, age`
	var createdPerson models.Person
	err = tx.QueryRowContext(ctx, query, person.FirstName, person.LastName, person.Type, person.Age).Scan(
		&createdPerson.ID,
		&createdPerson.FirstName,
		&createdPerson.LastName,
		&createdPerson.Type,
		&createdPerson.Age,
	)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("failed to insert person: %w", err)
	}

	if len(person.Courses) > 0 {
		insertCoursesQuery := `INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)`
		for _, courseID := range person.Courses {
			_, err := tx.ExecContext(ctx, insertCoursesQuery, createdPerson.ID, courseID)
			if err != nil {
				tx.Rollback()
				return models.Person{}, fmt.Errorf("failed to insert courses: %w", err)
			}
		}
	}

	selectCoursesQuery := `SELECT course_id FROM person_course WHERE person_id = $1`
	rows, err := tx.QueryContext(ctx, selectCoursesQuery, createdPerson.ID)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("failed to retrieve courses: %w", err)
	}
	defer rows.Close()

	var courses []int
	for rows.Next() {
		var courseID int
		if err := rows.Scan(&courseID); err != nil {
			tx.Rollback()
			return models.Person{}, fmt.Errorf("failed to scan course ID: %w", err)
		}
		courses = append(courses, courseID)
	}

	createdPerson.Courses = courses

	err = tx.Commit()
	if err != nil {
		return models.Person{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdPerson, nil
}

func (s *PersonService) DeletePerson(ctx context.Context, lastName string) error {
	tx, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	deleteCoursesQuery := `DELETE FROM person_course WHERE person_id IN 
                          (SELECT id FROM person WHERE LOWER(last_name) = LOWER($1))`
	_, err = tx.ExecContext(ctx, deleteCoursesQuery, lastName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete courses: %w", err)
	}

	deletePersonQuery := `DELETE FROM person WHERE LOWER(last_name) = LOWER($1)`
	result, err := tx.ExecContext(ctx, deletePersonQuery, lastName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no person found with last name: %s", lastName)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
