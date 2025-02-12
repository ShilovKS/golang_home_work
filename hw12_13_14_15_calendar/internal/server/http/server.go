package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Server struct {
	server *http.Server
	logger Logger
	app    Application
	config Config // Добавляем конфигурацию
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, config Config) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (s *Server) Start(_ context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handleHelloWorld)

	handler := loggingMiddleware(mux)

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
	}

	SetLogger(s.logger)

	s.logger.Info(fmt.Sprintf("Server started address: %s", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleHelloWorld(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
