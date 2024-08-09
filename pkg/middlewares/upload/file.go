package upload_middleware

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/blackmamoth/tasknet/pkg/utils"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) CheckFilePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			utils.SendAPIErrorResponse(w, http.StatusUnsupportedMediaType, fmt.Errorf("invalid content-type, should contain \"multipart/form-data\""))
			return
		}

		_, fileHeader, err := r.FormFile(config.GlobalConfig.AppConfig.FILE_OBJECT_NAME)

		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, "cannot process execution script, please make sure you've uploaded a file")
			return
		}

		ext := filepath.Ext(fileHeader.Filename)

		if ext != ".sh" {
			utils.SendAPIErrorResponse(w, http.StatusUnsupportedMediaType, "as of now tasknet only supports bash scripts execution, please only provide a bash script")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RenameFile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile(config.GlobalConfig.AppConfig.FILE_OBJECT_NAME)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		defer file.Close()

		newName, err := utils.GenerateRandomHex(8)

		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("an error occured while uploading your script"))
			return
		}

		handler.Header.Add("original_name", fmt.Sprintf(handler.Filename))

		handler.Filename = fmt.Sprintf("%s%s", newName, filepath.Ext(handler.Filename))

		next.ServeHTTP(w, r)

	})
}
