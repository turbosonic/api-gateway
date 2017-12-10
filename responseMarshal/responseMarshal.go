package responseMarshal

import (
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func AddHeaders(h http.Handler) http.Handler {
	addH := func(w http.ResponseWriter, r *http.Request) {
		u := uuid.NewV4()

		r.Header.Add("request_id", fmt.Sprint(u))

		w.Header().Set("request_id", fmt.Sprint(u))
		w.Header().Set("content-type", "application/json")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(addH)
}

func CorsHandler(h http.Handler) http.Handler {
	corsH := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type, origin, accept")
			w.WriteHeader(http.StatusOK)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(corsH)
}
