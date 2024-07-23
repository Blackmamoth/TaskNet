package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/db"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

var AccessTokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(config.GlobalConfig.AppConfig.ACCESS_TOKEN_SECRET), nil)

var RefreshTokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(config.GlobalConfig.AppConfig.REFRESH_TOKEN_SECRET), nil)

func SendAPIResponse(w http.ResponseWriter, status int, data any, cookies ...*http.Cookie) error {
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(generateAPIResponseBody(status, data))
}

func SendAPIErrorResponse(w http.ResponseWriter, status int, err interface{}) {
	if e, ok := err.(error); ok {
		SendAPIResponse(w, status, e.Error())
	} else {
		SendAPIResponse(w, status, err)
	}
}

func generateAPIResponseBody(status int, data any) map[string]any {
	if status >= 400 {
		return map[string]any{"status": status, "error": data}
	}
	return map[string]any{"status": status, "data": data}
}

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body should not be empty")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func generateMsgForField(fe validator.FieldError, v interface{}) (string, string) {
	t := reflect.TypeOf(v)

	field, _ := t.FieldByName(fe.StructField())

	jsonTag := field.Tag.Get("json")

	switch fe.Tag() {
	case "required":
		return jsonTag, fmt.Sprintf("\"%s\" is required", jsonTag)
	case "email":
		return jsonTag, fmt.Sprintf("\"%s\" must be a valid email address", jsonTag)
	case "min":
		return jsonTag, fmt.Sprintf("\"%s\" should contain at least %s characters", jsonTag, fe.Param())
	case "max":
		return jsonTag, fmt.Sprintf("\"%s\" should contain at most %s characters", jsonTag, fe.Param())
	case "dive":
		return jsonTag, fmt.Sprintf("\"%s\" should be in an array", jsonTag)
	}

	return fe.Field(), fe.Error()
}

func GenerateValidationErrorObject(ve validator.ValidationErrors, v interface{}) map[string]string {
	errs := map[string]string{}
	for _, fe := range ve {
		key, value := generateMsgForField(fe, v)
		errs[key] = value
	}
	return errs
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func calculateSecondsUntilEOD() (time.Time, time.Duration) {
	now := time.Now()

	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())

	duration := endOfDay.Sub(now)

	return now.Add(duration), duration
}

func SignAccessToken(r *http.Request, userId string) (string, error) {
	claims := map[string]interface{}{"id": userId, "remote_address": r.RemoteAddr}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Minute*time.Duration(config.GlobalConfig.AppConfig.ACCESS_TOKEN_EXPIRY_IN_MINS)))
	jwtauth.SetIssuedNow(claims)

	_, tokenString, err := AccessTokenAuth.Encode(claims)
	return tokenString, err
}

func SignRefreshToken(r *http.Request, userId string) (string, error) {
	claims := map[string]interface{}{"id": userId, "remote_address": r.RemoteAddr}

	eodInSeconds, eodInDuration := calculateSecondsUntilEOD()
	jwtauth.SetExpiry(claims, eodInSeconds)
	jwtauth.SetIssuedNow(claims)

	_, tokenString, err := RefreshTokenAuth.Encode(claims)

	db.RedisClient.Set(context.TODO(), userId, tokenString, eodInDuration)
	return tokenString, err
}

func GetRedisValue(key string) (string, error) {
	return db.RedisClient.Get(context.Background(), key).Result()
}
