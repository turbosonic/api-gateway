package logging

import (
	"net/http"
	"runtime"
	"time"
)

type LogClient interface {
	Log(*Log, string, string)
}

type LogHandler struct {
	client LogClient
}

type Log struct {
	Date          time.Time `json:"date"`
	RequestID     string    `json:"request_id"`
	Config        string    `json:"config"`
	Path          string    `json:"path"`
	URL           string    `json:"url"`
	Method        string    `json:"method"`
	StatusCode    int       `json:"status"`
	Duration      float64   `json:"duration"`
	ContentLength int64     `json:"content-length"`
	Host          string    `json:"host"`
	RemoteAddr    string    `json:"remote-address"`
	Agent         string    `json:"user-agent"`
	OS            string    `json:"os"`
	GoVersion     string    `json:"go-version"`
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
				r.Header.Get("config"),
				r.Header.Get("route"),
				r.RequestURI,
				r.Method,
				lrw.statusCode,
				float64(time.Since(start)) / float64(time.Millisecond),
				r.ContentLength,
				r.Host,
				r.RemoteAddr,
				r.Header.Get("User-Agent"),
				runtime.GOOS,
				runtime.Version()}
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
