package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/turbosonic/api-gateway/authentication/clients/auth0"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(h http.Handler) http.Handler {
	authProvider := os.Getenv("AUTHENTICATION_PROVIDER")

	switch authProvider {
	case "auth0":
		fmt.Println("[x] Authentication using Auth0")
		return auth0.CheckJwt(h)
	case "none":
		fmt.Println("[x] No authentication provider, all requests will be allowed")
		return noProvider(h)
	default:
		fmt.Println("[x] No authentication provider: all requests will be denied")
		return denyEverything(h)
	}
}

func noProvider(h http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(auth)
}

func denyEverything(h http.Handler) http.Handler {
	deny := func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Message: "No authentication provider has been provided",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	}
	return http.HandlerFunc(deny)
}
