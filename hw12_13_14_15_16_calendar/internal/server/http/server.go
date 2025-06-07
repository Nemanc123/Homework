package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	genApi "github.com/Calendar/hw12_13_14_15_calendar/api/gen/go"
	api "github.com/Calendar/hw12_13_14_15_calendar/internal/server/api"
	"github.com/Calendar/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msh string)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	GetEvent(ctx context.Context) ([]storage.Event, error)
	UpdateEvent(ctx context.Context, id int, event storage.Event) error
}

func NewServer(logger Logger, app Application) *Server {
	mux := http.NewServeMux()
	//mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	EventAPIService := api.NewEventAPIService(app)
	EventAPIController := genApi.NewEventAPIController(EventAPIService)

	router := genApi.NewRouter(EventAPIController)
	mux.Handle("/", router)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(mux, logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{server: server, logger: logger, app: app}
}

//	func homeHandler(w http.ResponseWriter, r *http.Request) {
//		if r.URL.Path != "/" {
//			http.NotFound(w, r)
//			return
//		}
//
// }
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start the server: %w", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to stop the server: %w", err)
	}
	return nil
}
