package logging

import (
	"net/http"
	"time"
)

type LogClient interface {
	Log(*Log, string, string)
}

type LogHandler struct {
	client LogClient
}

type Log struct {
	Date       time.Time `json:"date"`
	RequestID  string    `json:"request_id"`
	URL        string    `json:"url"`
	Method     string    `json:"method"`
	StatusCode int       `json:"status"`
	Duration   float64   `json:"duration"`
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func New(lc LogClient) LogHandler {
	return LogHandler{lc}
}

func (lh LogHandler) LogHandlerFunc(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)

		h.ServeHTTP(lrw, r)

		go func() {
			l := Log{
				start,
				r.Header.Get("request_id"),
				r.RequestURI,
				r.Method,
				lrw.statusCode,
				float64(time.Since(start)) / float64(time.Millisecond)}

			lh.client.Log(&l, "api-gateway", "request")
		}()
	}

	return http.HandlerFunc(auth)
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
