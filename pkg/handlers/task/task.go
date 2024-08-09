package task_handler

import (
	"fmt"
	"net/http"

	"github.com/blackmamoth/tasknet/pkg/config"
	upload_middleware "github.com/blackmamoth/tasknet/pkg/middlewares/upload"
	script_model "github.com/blackmamoth/tasknet/pkg/models/script"
	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/blackmamoth/tasknet/pkg/utils"
	"github.com/blackmamoth/tasknet/pkg/validations"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	taskService   types.TaskService
	userService   types.UserService
	scriptService types.ScriptService
}

func New(taskService types.TaskService, userService types.UserService, scriptService types.ScriptService) *Handler {
	return &Handler{
		taskService:   taskService,
		userService:   userService,
		scriptService: scriptService,
	}
}

func (h *Handler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(jwtauth.Verifier(utils.AccessTokenAuth))
	r.Use(jwtauth.Authenticator(utils.AccessTokenAuth))

	r.Post("/create", h.createTask)
	r.Group(func(r chi.Router) {
		uploadMiddleware := upload_middleware.New()
		r.Use(uploadMiddleware.CheckFilePayload)
		r.Use(uploadMiddleware.RenameFile)
		r.Post("/upload-execution-script", h.uploadExecutionScript)
	})
	r.Post("/get-tasks", h.getTasks)

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

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("an error occured while authorizing user claims"))
		return
	}

	err = h.taskService.CreateTask(types.CreateTaskSchema{
		Name:          payload.Name,
		ExecutionMode: payload.ExecutionMode,
		Priority:      payload.Priority,
		UserId:        claims["user_id"].(string),
	})

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendAPIResponse(w, http.StatusCreated, map[string]interface{}{"message": "Task created successfully."})
}

func (h *Handler) uploadExecutionScript(w http.ResponseWriter, r *http.Request) {

	if r.Form == nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Errorf("please provide all the required fields"))
		return
	}

	payload := validations.UploadExecutionScriptSchema{
		TaskId:     r.Form.Get("task_id"),
		Parameters: r.Form.Get("parameters"),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	task, err := h.taskService.GetTaskById(payload.TaskId)

	if err != nil || task == nil {
		utils.SendAPIErrorResponse(w, http.StatusNotFound, fmt.Errorf("task [%s] not found", payload.TaskId))
		return
	}

	if task.ScriptId.Valid {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Errorf("task [%s] already has a script associated to it", payload.TaskId))
		return
	}

	file, fileHeader, err := r.FormFile(config.GlobalConfig.AppConfig.FILE_OBJECT_NAME)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("an error occured while processing your file"))
		return
	}

	err = utils.SaveFile(file, fileHeader)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("an error occured while processing your file: %v", err))
		return
	}

	err = h.scriptService.CreateScript(&script_model.Script{
		Name:         fileHeader.Filename,
		OriginalName: fileHeader.Header.Get("original_name"),
		Parameters:   payload.Parameters,
	})

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("an error occured while processing your file: %v", err))
		return
	}

	s, err := h.scriptService.GetScriptByName(fileHeader.Filename)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("an error occured while processing your file: %v", err))
		return
	}

	err = h.taskService.RegisterScriptToTask(s.Id, task.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("an error occured while processing your file: %v", err))
		return
	}

	utils.SendAPIResponse(w, http.StatusOK, map[string]interface{}{"message": "You script was successfully uploaded."})
}

func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	var payload validations.GetTasksSchema

	utils.ParseJSON(r, &payload)

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("an error occured while authorizing user claims"))
		return
	}

	tasks, err := h.taskService.GetTasks(payload, claims["user_id"].(string))

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("an error occured while fetching your tasks: %v", err))
		return
	}

	data := map[string]interface{}{
		"tasks":   tasks,
		"message": "Tasks fetched successfully",
	}

	utils.SendAPIResponse(w, http.StatusOK, data)
}
