package initializer

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/turbosonic/api-gateway/logging"

	"github.com/turbosonic/api-gateway/parammap"

	"github.com/turbosonic/api-gateway/configurations"
	"github.com/turbosonic/api-gateway/relay"

	goji "goji.io"
	"goji.io/pat"
)

var (
	rel relay.Relay
)

func RegisterEndpoints(mux *goji.Mux, config configurations.Configuration, logger logging.LogClient) {
	rel = relay.New(10, 30, logger)

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
			//log.Println("Couldn't add method: %s", m)
			panic("Invalid method in configuration")
		}

		// move the method to a new variable
		d := m

		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			// check roles
			if !checkRoles(r, d) {
				w.WriteHeader(http.StatusNotFound)
				io.Copy(w, strings.NewReader("404 page not found"))
				return
			}

			// check scopes
			if !checkScopes(r, d) {
				w.WriteHeader(http.StatusNotFound)
				io.Copy(w, strings.NewReader("404 page not found"))
				return
			}

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
			request.Host = d.Destination.Host

			request.Header.Add("route", strings.Replace(p.String(), configName, "", 1))
			request.Header.Add("config", configName)

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

func checkRoles(r *http.Request, method configurations.EndpointMethod) bool {
	ur := r.Context().Value("roles").(string)
	for _, mr := range method.Roles {
		if mr == "*" || isThisInThat(mr, ur) {
			return true
		}
	}
	return false
}

func checkScopes(r *http.Request, method configurations.EndpointMethod) bool {
	us := r.Context().Value("scopes").(string)
	for _, ms := range method.Scopes {
		if ms == "*" || isThisInThat(ms, us) {
			return true
		}
	}
	return false
}

func isThisInThat(this string, that string) bool {
	for _, i := range strings.Split(that, " ") {
		if this == i {
			return true
		}
	}
	return false
}
