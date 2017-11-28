package initializer

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/turbosonic/api-gateway/parammap"

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

		ep := configName + endpoint.URL

		switch m.Method {
		case "GET":
			p = pat.Get(ep)
		case "POST":
			p = pat.Post(ep)
		case "PUT":
			p = pat.Put(ep)
		case "DELETE":
			p = pat.Delete(ep)
		default:
			log.Println("Couldn't add method: %s", m)
			panic("Invalid method in configuration")
		}

		// move the method to a new variable
		d := m

		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			// TODO: authorization

			// substitute parameters
			destinationURL := d.Destination.URL
			for p, v := range parammap.GetParams(ep, r) {
				destinationURL = strings.Replace(destinationURL, p, v, 1)
			}

			// append query string
			if r.URL.RawQuery != "" {
				destinationURL = destinationURL + "?" + r.URL.RawQuery
			}

			body, err := ioutil.ReadAll(r.Body)

			request := relay.RelayRequest{}
			request.URL = d.Destination.Host + destinationURL
			request.Method = d.Method
			request.Body = body
			request.Header = r.Header

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
