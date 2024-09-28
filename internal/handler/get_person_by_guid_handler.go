package handler

import (
	"net/http"

	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/services"
)

func GetPersonByGuid(service *services.PersonService) Route {
	return NewRoute("GET /api/person/{guid}", 
		func (w http.ResponseWriter, r *http.Request) {
			resp, err := service.GetOneByGuid(r.PathValue("guid"))
			encodeResponse(w, &applog.AppLogger{}, resp, err)
		},
	) 
}
