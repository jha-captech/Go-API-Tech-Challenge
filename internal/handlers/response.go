package handlers

import (
	"encoding/json"
	"go-api-tech-challenge/internal/models"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type outputCourse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type outputPerson struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
	Age       int    `json:"age"`
	Courses   []int  `json:"courses,omitempty"`
}

// mapOutput maps a models.Course struct to an outputCourse struct.
func mapOutputCourse(course models.Course) outputCourse {
	return outputCourse{
		ID:   course.ID,
		Name: course.Name,
	}
}

// mapMultipleOutput maps a slice of []models.Course to a slice of []outputCourse.
func mapMultipleOutputCourse(course []models.Course) []outputCourse {
	coursesOut := make([]outputCourse, len(course))
	for i := 0; i < len(course); i++ {
		courseOut := mapOutputCourse(course[i])
		coursesOut[i] = courseOut
	}

	return coursesOut
}

func mapOutputPerson(person models.Person) outputPerson {
	intCourseIDs := make([]int, len(person.Courses))
	for i, id := range person.Courses {
		intCourseIDs[i] = int(id) // Convert each int64 to int
	}

	return outputPerson{
		ID:        person.ID,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Type:      person.Type,
		Age:       person.Age,
		Courses:   intCourseIDs,
	}
}

func mapMultipleOutputPerson(person []models.Person) []outputPerson {
	personsOut := make([]outputPerson, len(person))
	for i := 0; i < len(person); i++ {
		personOut := mapOutputPerson(person[i])
		personsOut[i] = personOut
	}

	return personsOut

}

type responseCourse struct {
	Course outputCourse `json:"course"`
}

type responseCourses struct {
	Courses []outputCourse `json:"courses"`
}

type responseMsg struct {
	Message string `json:"message"`
}

type responsePerson struct {
	Person outputPerson `json:"person"`
}

type responsePersons struct {
	Persons []outputPerson `json:"persons"`
}

//type responseID struct {
//ObjectID int `json:"object_id"`
//}

type responseErr struct {
	Error            string    `json:"error,omitempty"`
	ValidationErrors []problem `json:"validation_errors,omitempty"`
}

// encodeResponse encodes data as a JSON response.
func encodeResponse(w http.ResponseWriter, logger *httplog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}
