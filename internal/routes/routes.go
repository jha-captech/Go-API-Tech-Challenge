package routes

import (
	"go-api-tech-challenge/internal/handlers"
	"go-api-tech-challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type Option func(*routerOptions)

type routerOptions struct {
	registerHealthRoute bool
}

// WithRegisterHealthRoute controls whether a healthcheck route will be registered. If `false` is
// passed in or this function is not called, the default is `false`.
func WithRegisterHealthRoute(registerHealthRoute bool) Option {
	return func(options *routerOptions) {
		options.registerHealthRoute = registerHealthRoute
	}
}

func RegisterRoutes(router *chi.Mux, logger *httplog.Logger, svsCourse *services.CourseService, svsPerson *services.PersonService, opts ...Option) {

	options := routerOptions{
		registerHealthRoute: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	router.Route("/api", func(router chi.Router) {

		if options.registerHealthRoute {
			router.Get("/health-check", handlers.HandleHealth(logger))
		}

		router.Route("/course", func(router chi.Router) {

			router.Get("/", handlers.HandleListCourses(logger, svsCourse))
			router.Post("/", handlers.HandleCreateCourse(logger, svsCourse))
			router.Get("/{ID}", handlers.HandleGetCourseByID(logger, svsCourse))
			router.Put("/{ID}", handlers.HandleUpdateCourse(logger, svsCourse))
			router.Delete("/{ID}", handlers.HandleDeleteCourse(logger, svsCourse))

		})
		router.Route("/person", func(router chi.Router) {

			router.Get("/", handlers.HandleListPersons(logger, svsPerson))
			router.Post("/", handlers.HandleCreatePerson(logger, svsPerson))
			router.Get("/{name}", handlers.HandleGetPersonByName(logger, svsPerson))
			router.Put("/{name}", handlers.HandleUpdatePerson(logger, svsPerson))
			router.Delete("/{name}", handlers.HandleDeletePerson(logger, svsPerson))

		})
	})

}
