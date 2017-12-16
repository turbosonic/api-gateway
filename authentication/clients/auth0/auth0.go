package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	auth0 "github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
)

type Response struct {
	Message string `json:"message"`
}

func CheckJwt(h http.Handler) http.Handler {
	fmt.Println("[x] Auth0 domain:", os.Getenv("AUTH0_DOMAIN"))
	fmt.Println("[x] Auth0 audience:", os.Getenv("AUTH0_AUDIENCE"))

	JWKS_URI := "https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json"
	AUTH0_API_ISSUER := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"

	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: JWKS_URI})
	aud := os.Getenv("AUTH0_AUDIENCE")
	audience := []string{aud}
	configuration := auth0.NewConfiguration(client, audience, AUTH0_API_ISSUER, jose.RS256)
	validator := auth0.NewValidator(configuration)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println("Token is not valid or missing token")
			fmt.Println(err)

			response := Response{
				Message: "Missing or invalid token.",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)

		} else {
			claims := map[string]interface{}{}
			err := validator.Claims(r, token, &claims)
			if err != nil {
				response := Response{
					Message: "Missing or invalid token.",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
			} else {
				ctx := context.WithValue(r.Context(), "roles", claims[os.Getenv("AUTH0_AUDIENCE")+"/roles"])
				ctx = context.WithValue(ctx, "scopes", claims["scope"])
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
