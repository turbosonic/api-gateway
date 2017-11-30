package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/turbosonic/api-gateway/authentication"
	"github.com/turbosonic/api-gateway/configurations"
	"github.com/turbosonic/api-gateway/initializer"
	"github.com/turbosonic/api-gateway/logging"
	"github.com/turbosonic/api-gateway/logging/clients/elk"
	"github.com/turbosonic/api-gateway/logging/clients/stdout"
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
		*configFile = "config.yaml"
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

	// add response logging
	var logClient logging.LogClient

	elasticURL := os.Getenv("LOGGING_ELASTIC_URL")
	if elasticURL != "" {
		logClient = elk.New(elasticURL)
	} else {
		logClient = stdout.New()
	}

	logHandler := logging.New(logClient)
	mux.Use(logHandler.LogHandlerFunc)

	// Register the endpoints
	initializer.RegisterEndpoints(mux, config, logClient)

	// start listening
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err)
	}
}
