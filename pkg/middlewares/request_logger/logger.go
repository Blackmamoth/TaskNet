package logger_middleware

import (
	"fmt"
	"net/http"
	"time"

	golog "github.com/blackmamoth/GoLog"
	"github.com/blackmamoth/tasknet/pkg/config"
)

type RequestDetails struct {
	Method     string
	URL        string
	Proto      string
	RemoteAddr string
	StatusCode int
	Size       int
	Duration   time.Duration
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

var logger golog.Logger

func init() {
	logger = golog.New()
	if config.GlobalConfig.AppConfig.ENVIRONMENT == "DEVELOPMENT" {
		logger.Set_Log_Level(golog.LOG_LEVEL_DEBUG)
	}
	logger.Set_Log_Stream(golog.LOG_STREAM_FILE)
	logger.Set_File_Name(fmt.Sprintf("%s/%s", config.GlobalConfig.AppConfig.APP_LOG_PATH, config.GlobalConfig.AppConfig.APP_LOG_FILE))
	logger.With_Emoji(true)
	logger.Set_Log_Format("[%(asctime)] %(levelname) - %(message)")
	logger.Exit_On_Critical(true)
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

func HttpRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		var hostUrl string

		if r.TLS == nil {
			hostUrl = fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
		} else {
			hostUrl = fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
		}

		details := RequestDetails{
			Method:     r.Method,
			URL:        hostUrl,
			Proto:      r.Proto,
			RemoteAddr: r.RemoteAddr,
			StatusCode: lrw.statusCode,
			Size:       lrw.size,
			Duration:   duration,
		}

		logger.INFO("%s %s %s from %s - %d %dB in %s", details.Method, details.URL, details.Proto, details.RemoteAddr, details.StatusCode, details.Size, details.Duration)

	})
}
