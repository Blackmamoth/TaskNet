package utils

import (
	"encoding/json"
	"net/http"
)

func SendAPIResponse(w http.ResponseWriter, status int, data any, cookie *http.Cookie) error {
	if cookie != nil {
		http.SetCookie(w, cookie)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(generateAPIResponseBody(status, data))
}

func SendAPIErrorResponse(w http.ResponseWriter, status int, err error) {
	SendAPIResponse(w, status, err.Error(), nil)
}

func generateAPIResponseBody(status int, data any) map[string]any {
	if status >= 400 {
		return map[string]any{"status": status, "error": data}
	}
	return map[string]any{"status": status, "data": data}
}
