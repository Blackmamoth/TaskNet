package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

type APIServer struct {
	host string
	addr string
	conn *pgx.Conn
}

func NewAPIServer(host, addr string, conn *pgx.Conn) *APIServer {
	return &APIServer{
		host: host,
		addr: addr,
		conn: conn,
	}
}

func (s *APIServer) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Compress(5, "gzip"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIResponse(w, http.StatusOK, "Welcome to TaskNet the distributed task scheduler.", nil)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIErrorResponse(w, http.StatusNotFound, fmt.Errorf("route not found for [%s] %s", r.Method, r.URL.Path))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIErrorResponse(w, http.StatusMethodNotAllowed, fmt.Errorf("method [%s] not allowed for route %s", r.Method, r.URL.Path))
	})

	config.Logger.INFO("Application running on port: %s", s.addr)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.addr), r)
}
