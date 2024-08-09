package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackmamoth/tasknet/pkg/config"
	auth_handler "github.com/blackmamoth/tasknet/pkg/handlers/auth"
	task_handler "github.com/blackmamoth/tasknet/pkg/handlers/task"
	logger_middleware "github.com/blackmamoth/tasknet/pkg/middlewares/request_logger"
	script_repository "github.com/blackmamoth/tasknet/pkg/repository/script"
	task_repository "github.com/blackmamoth/tasknet/pkg/repository/task"
	user_repository "github.com/blackmamoth/tasknet/pkg/repository/user"
	script_service "github.com/blackmamoth/tasknet/pkg/services/script"
	task_service "github.com/blackmamoth/tasknet/pkg/services/task"
	user_service "github.com/blackmamoth/tasknet/pkg/services/user"
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
	r.Use(logger_middleware.HttpRequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Compress(5, "gzip"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.GlobalConfig.AppConfig.APP_FRONTEND},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIResponse(w, http.StatusOK, map[string]interface{}{"message": "Welcome to TaskNet the distributed task scheduler."})
	})

	r.Mount("/v1/api", s.registerRoutes())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIErrorResponse(w, http.StatusNotFound, fmt.Sprintf("route not found for [%s] %s", r.Method, r.URL.Path))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.SendAPIErrorResponse(w, http.StatusMethodNotAllowed, fmt.Sprintf("method [%s] not allowed for route %s", r.Method, r.URL.Path))
	})

	config.Logger.INFO("Application running on port: %s", s.addr)

	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.addr), r)
}

func (s *APIServer) registerRoutes() *chi.Mux {
	subRouter := chi.NewRouter()

	userRepository := user_repository.New(s.conn)
	userService := user_service.New(userRepository)
	authHandler := auth_handler.New(userService)

	scriptRepository := script_repository.New(s.conn)
	scriptService := script_service.New(scriptRepository)
	taskRepository := task_repository.New(s.conn)
	taskService := task_service.New(taskRepository)
	taskHandler := task_handler.New(taskService, userService, scriptService)

	subRouter.Mount("/auth", authHandler.RegisterRoutes())
	subRouter.Mount("/task", taskHandler.RegisterRoutes())

	return subRouter
}
