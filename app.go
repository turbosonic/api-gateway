package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/turbosonic/api-gateway/authentication"
	"github.com/turbosonic/api-gateway/configurations"
	"github.com/turbosonic/api-gateway/factories"
	"github.com/turbosonic/api-gateway/initializer"
	"github.com/turbosonic/api-gateway/logging"
	"github.com/turbosonic/api-gateway/responseMarshal"

	"github.com/zenazn/goji/graceful"
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

	// add response marshaling
	mux.Use(responseMarshal.CorsHandler)
	mux.Use(responseMarshal.AddHeaders)

	// add response logging
	logClient := factories.LogClient()
	logHandler := logging.New(logClient)
	mux.Use(logHandler.LogHandlerFunc)

	// add authentication
	mux.Use(authentication.Handler)

	// Register the endpoints
	initializer.RegisterEndpoints(mux, config, logClient)

	// get the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("[x] Listening on port", port)

	// check for cert and key
	if _, certErr := os.Stat("./data/certs/cert.pem"); os.IsNotExist(certErr) {
		fmt.Println("[x] Not using TLS")
		err = http.ListenAndServe(":"+port, mux)
		if err != nil {
			log.Println(err)
		}
	} else {
		fmt.Println("[x] Using TLS")
		err = graceful.ListenAndServeTLS(":"+port, "./data/certs/cert.pem", "./data/certs/key.pem", mux)
		if err != nil {
			log.Println(err)
		}
	}

}
