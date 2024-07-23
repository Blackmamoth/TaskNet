package task_handler

import (
	"fmt"
	"net/http"

	auth_middleware "github.com/blackmamoth/tasknet/pkg/middlewares/auth"
	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/blackmamoth/tasknet/pkg/utils"
	"github.com/blackmamoth/tasknet/pkg/validations"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	taskService types.TaskService
	userService types.UserService
}

func New(taskService types.TaskService, userService types.UserService) *Handler {
	return &Handler{
		taskService: taskService,
		userService: userService,
	}
}

func (h *Handler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	authMiddleware := auth_middleware.New(h.userService)

	r.Use(jwtauth.Verifier(utils.AccessTokenAuth))
	r.Use(jwtauth.Authenticator(utils.AccessTokenAuth))
	r.Use(authMiddleware.VerifyRefreshToken)

	r.Post("/create", h.createTask)

	return r
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var payload validations.CreateTaskSchema

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Errorf("please provide all the required fields"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	_, err := h.taskService.GetTaskByName(payload.Name)

	if err == nil {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Errorf("task with name [%s] already exists", payload.Name))
		return
	}

	err = h.taskService.CreateTask(payload)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendAPIResponse(w, http.StatusCreated, "Task created successfully")
}
