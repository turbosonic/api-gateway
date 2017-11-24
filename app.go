package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/turbosonic/api-gateway/authentication"
	"github.com/turbosonic/api-gateway/configurations"
	"github.com/turbosonic/api-gateway/initializer"
	"github.com/turbosonic/api-gateway/responseMarshal"

	goji "goji.io"
)

var (
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

	// add response marshaling
	mux.Use(responseMarshal.AddHeaders)

	// Register the endpoints
	initializer.RegisterEndpoints(mux, config)

	// start listening
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err)
	}
}
