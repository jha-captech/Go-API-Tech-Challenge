package handler

import (
	"net/http"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/models"
	"jf.go.techchallenge/internal/services"
)

func GetAllPersons(service *services.Person, logger *applog.AppLogger) Route {
	return NewRoute("GET /api/person", func(w http.ResponseWriter, r *http.Request) {
		persons, err := service.GetAll(r.URL.Query())
		encodeResponse(w, logger, persons, err)
	})
}

func CreatePerson(service *services.Person, logger *applog.AppLogger) Route {
	return NewRoute("POST /api/person", func(w http.ResponseWriter, r *http.Request) {

		input, err := decodeBody[models.PersonInput](r)

		if err != nil {
			logger.Debug("Failed to Deserialize Request Body ", err)
			encodeError(w, apperror.BadRequest("Invalid JSON"))
			return
		}

		person, err := service.Create(*input)
		encodeCreated(w, logger, person, err)
	})
}

func GetOnePerson(service *services.Person, logger *applog.AppLogger) Route {
	return NewRoute("GET /api/person/{guid}", func(w http.ResponseWriter, r *http.Request) {
		person, err := service.GetOneByGuid(r.PathValue("guid"))
		encodeResponse(w, logger, person, err)
	})
}

func UpdateOnePerson(service *services.Person, logger *applog.AppLogger) Route {
	return NewRoute("PUT /api/person/{guid}", func(w http.ResponseWriter, r *http.Request) {
		input, err := decodeBody[models.PersonInput](r)

		if err != nil {
			logger.Debug("Failed to Deserialize Request Body ", err)
			encodeError(w, apperror.BadRequest("Invalid JSON"))
			return
		}

		person, err := service.Update(r.PathValue("guid"), *input)
		encodeResponse(w, logger, person, err)
	})
}

func DeleteOnePerson(service *services.Person, logger *applog.AppLogger) Route {
	return NewRoute("DELETE /api/person/{guid}", func(w http.ResponseWriter, r *http.Request) {
		err := service.Delete(r.PathValue("guid"))
		encodeResponse(w, logger, "OK", err)
	})
}
