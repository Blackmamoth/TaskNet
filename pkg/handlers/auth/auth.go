package auth_handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/db"
	auth_middleware "github.com/blackmamoth/tasknet/pkg/middlewares/auth"
	user_model "github.com/blackmamoth/tasknet/pkg/models/user"
	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/blackmamoth/tasknet/pkg/utils"
	"github.com/blackmamoth/tasknet/pkg/validations"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(utils.AccessTokenAuth))
		r.Use(jwtauth.Authenticator(utils.AccessTokenAuth))
		r.Post("/logout", h.logout)
	})

	r.Group(func(r chi.Router) {
		authMiddleware := auth_middleware.New(h.userService)
		r.Use(authMiddleware.VerifyRefreshToken)
		r.Post("/refresh", h.refresh)
	})

	return r
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload validations.RegisterUserSchema

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("please provide all the required fields"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	_, err := h.userService.GetUserByUsername(payload.Username)

	if err == nil {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Sprintf("username [%s] already in use", payload.Username))
		return
	}

	_, err = h.userService.GetUserByEmail(payload.Email)

	if err == nil {
		utils.SendAPIErrorResponse(w, http.StatusConflict, fmt.Sprintf("email [%s] already in use", payload.Email))
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("an error occured while processing your password"))
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

	utils.SendAPIResponse(w, http.StatusCreated, map[string]interface{}{"message": "Your registration was successful"})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload validations.LoginUserSchema

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("please provide all the required fields"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := utils.GenerateValidationErrorObject(err.(validator.ValidationErrors), payload)
		utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, errs)
		return
	}

	u, err := h.userService.GetUserByUsername(payload.Username)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("invalid username please check your username and try again"))
		return
	}

	if !utils.ComparePassword(payload.Password, u.Password) {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("invalid password please check your password and try again"))
		return
	}

	accessToken, err := utils.SignAccessToken(r, u.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("an error occured while processing your credentials"))
		return
	}

	refreshToken, err := utils.SignRefreshToken(r, u.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("an error occured while processing your credentials"))
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
		"message": "Login successful",
	}

	now := time.Now()

	refreshTokenCookie := utils.GenerateRefreshTokenCookie(refreshToken, time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location()))

	utils.SendAPIResponse(w, http.StatusOK, data, &refreshTokenCookie)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, "an error occured while authenticating your request")
		return
	}

	userId := claims["user_id"].(string)

	u, err := h.userService.GetUserById(userId)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusNotFound, fmt.Sprintf("user [%s] not found", err))
		return
	}

	db.RedisClient.Del(context.TODO(), u.Id)

	refreshTokenCookie := utils.GenerateRefreshTokenCookie("", time.Unix(0, 0))

	utils.SendAPIResponse(w, http.StatusOK, map[string]interface{}{"message": "logout successful"}, &refreshTokenCookie)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var refreshToken string
	if len(r.Cookies()) > 0 {
		for _, cookie := range r.Cookies() {
			if cookie.Name == config.GlobalConfig.AppConfig.REFRESH_TOKEN_NAME {
				refreshToken = cookie.Value
				break
			}
		}
	}

	if refreshToken == "" {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, "refresh token cookie not found")
		return
	}

	token, err := jwtauth.VerifyToken(utils.RefreshTokenAuth, refreshToken)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, "could not verify jwt token")
		return
	}

	user_id, exists := token.Get("user_id")

	if !exists {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, "could not verify jwt token")
		return
	}

	u, err := h.userService.GetUserById(user_id.(string))

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, "invalid user")
		return
	}

	accessToken, err := utils.SignAccessToken(r, u.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("an error occured while processing your credentials"))
		return
	}

	refreshToken, err = utils.SignRefreshToken(r, u.Id)

	if err != nil {
		utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("an error occured while processing your credentials"))
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
		"message": "Refresh successful",
	}

	now := time.Now()

	refreshTokenCookie := utils.GenerateRefreshTokenCookie(refreshToken, time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location()))

	utils.SendAPIResponse(w, http.StatusOK, data, &refreshTokenCookie)

}
