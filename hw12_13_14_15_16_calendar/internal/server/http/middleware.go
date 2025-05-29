package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	Status int
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func loggingMiddleware(next http.Handler, logger Logger) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		Response := &ResponseWriterWrapper{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(Response, r)
		logMsg := fmt.Sprintf("%s[%v] %s %s %s %d %dms",
			r.RemoteAddr,
			start,
			r.Method,
			r.URL.Path,
			r.Proto,
			Response.Status,
			time.Since(start).Milliseconds(),
		)
		if 200 <= Response.Status || Response.Status < 300 {
			logger.Info(logMsg)
		} else {
			logger.Error(logMsg)
		}
	})
}

//66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
