package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jf.go.techchallenge/internal/config"
	"jf.go.techchallenge/internal/handler"
)

func main() {
	fx.New(
		fx.Provide(
			handler.NewGetPersonsHanlder,
			NewServeMux,
			NewHTTPServer,
			// services.NewPersonService,
			// NewDatabase,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func NewServeMux(route *handler.GetPersonsHandler) *http.ServeMux {
	mux := http.NewServeMux()
	// for _, route := range routes {
	mux.Handle(route.Pattern(), route)
	// }
	return mux
}

// Start Database
func NewDatabase() (*gorm.DB, error) {
	config, err := config.New()

	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)

	// Connect to database.
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{
		// Logger: newLogger, todo
	})
}

// Start Http Server.
func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
