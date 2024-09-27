package handler

import (
	"net/http"

	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/services"
)

type GetPersonsHandler struct {
	service *services.PersonService
}

func NewGetPersonsHanlder(service *services.PersonService) Route {
	return &GetPersonsHandler{
		service: service,
	}
}

func (*GetPersonsHandler) Pattern() string {
	return "GET /api/person"
}

var validFilters = services.MakeFilterColumns(services.ValidFilters{
	"FirstName",
	"LastName",
	"Email",
})

func (s *GetPersonsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	filters, err := services.ParseURLFilters(r.URL.Query(), validFilters)

	if err != nil {
		return
	}
	persons, _ := s.service.GetPersons(filters)

	encodeResponse(w, &applog.AppLogger{}, 200, persons)
}
