package auth_handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackmamoth/tasknet/pkg/config"
	user_model "github.com/blackmamoth/tasknet/pkg/models/user"
	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/blackmamoth/tasknet/pkg/utils"
	"github.com/blackmamoth/tasknet/pkg/validations"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	userService types.UserService
}

func New(userService types.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", h.registerUser)
	r.Post("/login", h.login)

	return r
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload validations.RegisterUserSchema

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Errorf("please provide all the required fields"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	_, err := h.userService.GetUserByUsername(payload.Username)

	if err == nil {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Errorf("username [%s] already in use", payload.Username))
		return
	}

	_, err = h.userService.GetUserByEmail(payload.Email)

	if err == nil {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Errorf("email [%s] already in use", payload.Email))
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("an error occured while processing your password"))
		return
	}

	err = h.userService.CreateUser(&user_model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendAPIResponse(w, http.StatusCreated, "Your registration was successful")
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload validations.LoginUserSchema

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Errorf("please provide all the required fields"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	u, err := h.userService.GetUserByUsername(payload.Username)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid username please check your username and try again"))
		return
	}

	if !utils.ComparePassword(payload.Password, u.Password) {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid password please check your password and try again"))
		return
	}

	accessToken, err := utils.SignAccessToken(r, u.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("an error occured while processing your credentials"))
		return
	}

	refreshToken, err := utils.SignRefreshToken(r, u.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("an error occured while processing your credentials"))
		return
	}

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"user_id":    u.Id,
			"username":   u.Username,
			"email":      u.Email,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		},
		config.GlobalConfig.AppConfig.ACCESS_TOKEN_NAME: accessToken,
	}

	now := time.Now()

	refreshTokenCookie := http.Cookie{
		Name:     config.GlobalConfig.AppConfig.REFRESH_TOKEN_NAME,
		Value:    refreshToken,
		Secure:   config.GlobalConfig.AppConfig.ENVIRONMENT != "PRODUCTION",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location()),
	}

	utils.SendAPIResponse(w, http.StatusOK, data, &refreshTokenCookie)
}
