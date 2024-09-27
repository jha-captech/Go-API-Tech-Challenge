package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type GetPersonsHandler struct{}

func NewGetPersonsHanlder() *GetPersonsHandler {
	return &GetPersonsHandler{}
}

func (*GetPersonsHandler) Pattern() string {
	return "POST /api/person"
}

func (*GetPersonsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
	}
}
