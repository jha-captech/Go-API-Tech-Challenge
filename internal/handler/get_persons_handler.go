package handler

import (
	"fmt"
	"net/http"

	"jf.go.techchallenge/internal/applog"
	"jf.go.techchallenge/internal/services"
)

var validFilters = services.MakeFilterColumns(services.ValidFilters{
	"FirstName",
	"LastName",
	"Email",
})

func NewGetPersonsHanlder(service *services.PersonService) Route {
	return NewRoute("GET /api/person", func(w http.ResponseWriter, r *http.Request) {
		filters, err := services.ParseURLFilters(r.URL.Query(), validFilters)

		if err != nil {
			fmt.Println("PRinting errors ", err)
			encodeError(w, err)
			return
		}
		persons, err := service.GetPersons(filters)

		encodeResponse(w, &applog.AppLogger{}, persons, err)
	})
}
