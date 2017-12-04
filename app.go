package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/turbosonic/api-gateway/authentication/clients/auth0"
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
		*configFile = "config.yaml"
	}

	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	// get all of the endpoints
	config, err := configurations.GetConfiguration(*configFile)
	if err != nil {
		panic("Could not read configurations")
	}

	// create a new mux from goju
	mux := goji.NewMux()

	// add authentication
	mux.Use(auth0.CheckJwt)

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
