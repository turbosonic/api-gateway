package logging

import (
	"net/http"
	"time"

	elastic "gopkg.in/olivere/elastic.v2"
)

type LogHandler struct {
	client *elastic.Client
}

type log struct {
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

func New() LogHandler {
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}

	return LogHandler{client}
}

func (lh LogHandler) LogHandlerFunc(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)

		h.ServeHTTP(lrw, r)

		go func() {
			l := log{
				start,
				r.Header.Get("request_id"),
				r.RequestURI,
				r.Method,
				lrw.statusCode,
				float64(time.Since(start)) / float64(time.Millisecond)}

			lh.client.Index().Index("turbosonic").Type("log").Id(l.RequestID).BodyJson(l).Do()
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
