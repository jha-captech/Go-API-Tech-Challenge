package handler

import (
	"net/http"

	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/services"
)

type GetPersonByGuidHandler struct {
	service *services.PersonService
}

func NewGetPersonByGuid(service *services.PersonService) Route {
	return &GetPersonByGuidHandler{
		service: service,
	}
}

func (*GetPersonByGuidHandler) Pattern() string {
	return "GET /api/person/{guid}"
}

func (s *GetPersonByGuidHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	resp, err := s.service.GetOneByGuid(r.PathValue("guid"))
	encodeResponse(w, &applog.AppLogger{}, resp, err)
}
