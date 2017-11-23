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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(addH)
}
