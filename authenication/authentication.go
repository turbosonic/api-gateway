package authentication

import (
	"net/http"
)

func Authenticate(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(auth)
}
