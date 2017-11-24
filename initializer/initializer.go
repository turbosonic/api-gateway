package initializer

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/turbosonic/api-gateway/configurations"
	"github.com/turbosonic/api-gateway/relay"
	goji "goji.io"
	"goji.io/pat"
)

var (
	rel = relay.New(10, 30)
)

func RegisterEndpoints(mux *goji.Mux, config configurations.Configuration) {
	// loop through each endpoint and add a request handler
	for _, e := range config.Endpoints {
		createEndpoint(mux, config.Name, &e)
	}
}

func createEndpoint(mux *goji.Mux, configName string, endpoint *configurations.Endpoint) {
	for _, m := range endpoint.Methods {
		var p *pat.Pattern

		switch m.Method {
		case "GET":
			p = pat.Get(configName + endpoint.URL)
		case "POST":
			p = pat.Post(configName + endpoint.URL)
		case "PUT":
			p = pat.Put(configName + endpoint.URL)
		case "DELETE":
			p = pat.Delete(configName + endpoint.URL)
		default:
			log.Println("Couldn't add method: %s", m)
			panic("Invalid method in configuration")
		}

		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			// TODO: authorization

			// append query string
			destinationURL := m.Destination.URL
			if r.URL.RawQuery != "" {
				destinationURL = destinationURL + "?" + r.URL.RawQuery
			}

			request := relay.RelayRequest{}
			request.URL = m.Destination.Host + destinationURL
			request.Method = m.Method

			resp, err := rel.MakeRequest(request)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				io.Copy(w, strings.NewReader("404 page not found"))
				return
			}

			defer resp.Body.Close()
			w.WriteHeader(resp.StatusCode)

			io.Copy(w, resp.Body)
		})
	}
}
