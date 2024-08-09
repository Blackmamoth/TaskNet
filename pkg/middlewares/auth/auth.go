package auth_middleware

import (
	"fmt"
	"net/http"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/types"
	"github.com/blackmamoth/tasknet/pkg/utils"
	"github.com/go-chi/jwtauth/v5"
)

type Middleware struct {
	service types.UserService
}

func New(service types.UserService) *Middleware {
	return &Middleware{
		service: service,
	}
}

func (m *Middleware) VerifyRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		if len(cookies) == 0 {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("refresh token not present in cookies"))
			return
		}

		var refreshToken string

		for _, cookie := range cookies {
			if cookie.Name == config.GlobalConfig.AppConfig.REFRESH_TOKEN_NAME {
				refreshToken = cookie.Value
				break
			}
		}

		if refreshToken == "" {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("refresh token not present in cookies"))
			return
		}

		token, err := jwtauth.VerifyToken(utils.RefreshTokenAuth, refreshToken)

		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token"))
			return
		}

		userId, ok := token.Get("user_id")
		if !ok {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token"))
			return
		}

		_, err = m.service.GetUserById(userId.(string))

		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token, user not found"))
			return
		}

		_, err = utils.GetRedisValue(userId.(string))

		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token, expired"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
