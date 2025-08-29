package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggingDeps struct {
	Logger *logrus.Logger
}

func Logging(deps *LoggingDeps) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapper := &WrapperWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapper, r)

			fields := logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": wrapper.statusCode,
				"time":   time.Since(start),
			}

			if wrapper.statusCode >= 400 {
				message := fmt.Sprintf("%s %s %d %s", r.Method, r.URL.Path, wrapper.statusCode, time.Since(start))
				deps.Logger.WithFields(fields).Warn(message)
			} else if wrapper.statusCode >= 500 {
				message := fmt.Sprintf("%s %s %d %s", r.Method, r.URL.Path, wrapper.statusCode, time.Since(start))
				deps.Logger.WithFields(fields).Error(message)
			} else {
				message := fmt.Sprintf("%s %s %d %s", r.Method, r.URL.Path, wrapper.statusCode, time.Since(start))
				deps.Logger.WithFields(fields).Info(message)
			}
		})
	}
}
