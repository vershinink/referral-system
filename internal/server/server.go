// Пакет server используется для инициализации http сервера и
// работы с ним.
package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"referral-rest-api/internal/config"
	"referral-rest-api/internal/server/api/codes"
	"referral-rest-api/internal/server/api/users"
	"referral-rest-api/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

// gracefulStopTime определяет время на завершение всех операций сервера.
const gracefulStopTime time.Duration = 10 * time.Second

// Server - структура сервера.
type Server struct {
	srv *http.Server
	mux *chi.Mux
}

// New - конструктор сервера.
func New(cfg *config.Config) *Server {
	r := chi.NewRouter()
	server := &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		mux: r,
	}
	return server
}

// Start запускает HTTP сервер в отдельной горутине.
func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Print("failed to start server")
		}
	}()
}

// API инициализирует все обработчики API.
func (s *Server) API(app *service.App) {
	// Публичные эндпоинты.
	s.mux.Post("/api/users", users.Register(app))
	s.mux.Post("/api/users/login", users.Login(app))
	s.mux.Get("/api/users/{id}", users.UsersByReferrer(app))
	s.mux.Get("/api/refcodes/email/{email}", codes.CodeByEmail(app))

	// Защищенные эндпоинты.
	s.mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.JWTauth()))
		r.Use(jwtauth.Authenticator(app.JWTauth()))

		r.Post("/api/refcodes", codes.Create(app))
		r.Get("/api/refcodes", codes.CodeByID(app))
		r.Delete("/api/refcodes", codes.Delete(app))
	})
}

// Middleware инициализирует все обработчики middleware.
func (s *Server) Middleware() {
	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)

	opts := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	s.mux.Use(cors.Handler(opts))

}

// Shutdown останавливает сервер используя graceful shutdown
func (s *Server) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulStopTime)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}
}
