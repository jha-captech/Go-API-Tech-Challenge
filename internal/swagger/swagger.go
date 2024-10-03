package swagger

import (
	"fmt"
	"go-api-tech-challenge/internal/swagger/docs"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RunSwagger(r *chi.Mux, logger *httplog.Logger, host string) {
	// docs
	docs.SwaggerInfo.Title = "Go API Tech Challenge"
	docs.SwaggerInfo.Description = "Microservice for tech challenge"
	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = ""

	docs.SwaggerInfo.Schemes = []string{"http"}

	// handler
	baseURL := "http://" + host

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(baseURL+"/swagger/doc.json"),
	))

	log.Printf("Swagger URL: %s/swagger/index.html", baseURL)
	logger.Info(fmt.Sprintf("Swagger URL: %s/swagger/index.html", baseURL))
}
