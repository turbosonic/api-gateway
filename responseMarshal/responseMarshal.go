package responseMarshal

import (
	"fmt"
	"net/http"
	"os"

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

			acqo := os.Getenv("CORS_ORIGIN")
			if len(acqo) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", acqo) // no default
			}

			acma := os.Getenv("CORS_MAX_AGE")
			if len(acma) == 0 {
				acma = "600" // ten minute default
			}

			w.Header().Set("Access-Control-Max-Age", acma)
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
