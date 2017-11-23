package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/turbosonic/golddust/authenication"
	"github.com/turbosonic/golddust/configurations"
	"github.com/turbosonic/golddust/relay"
	"github.com/turbosonic/golddust/responseMarshal"

	goji "goji.io"
	"goji.io/pat"
)

var (
	rel        = relay.New(10, 30)
	configFile = flag.String("config", "", "yaml file for configuration")
)

func main() {
	// pick up the flags
	flag.Parse()

	if *configFile == "" {
		panic("No config file provided")
	}

	// get all of the endpoints
	config, err := configurations.GetConfiguration(*configFile)
	if err != nil {
		panic("Could not read configurations")
	}

	// create a new mux from goju
	mux := goji.NewMux()

	// add authentication
	mux.Use(authentication.Authenticate)

	// add response marshalling
	mux.Use(responseMarshal.AddHeaders)

	// loop through each endpoint and add a request handler
	for _, e := range config.Endpoints {
		createEndpoint(mux, config.Name, &e)
	}

	// start listening
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err)
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
			// authorization

			request := relay.RelayRequest{}
			request.URL = m.Destination.Host + m.Destination.URL
			request.Method = m.Method

			resp, err := rel.MakeRequest(request)
			if err != nil {
				return
			}

			defer resp.Body.Close()
			w.WriteHeader(resp.StatusCode)

			io.Copy(w, resp.Body)
		})
	}
}
