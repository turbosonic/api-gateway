package configurations

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Name      string
	Endpoints []Endpoint `yaml:",flow"`
}

type Endpoint struct {
	URL     string
	Methods []EndpointMethod `yaml:",flow"`
}

type EndpointMethod struct {
	Method      string
	Roles       []string `yaml:",flow"`
	Destination EndpointDestination
	//Destinations []EndpointDestination `yaml:"destinations"`
}

type EndpointDestination struct {
	Name string
	Host string
	URL  string
}

func GetConfiguration(filePath string) (Configuration, error) {

	filename, _ := filepath.Abs(filePath)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Could not read file")
		panic("Could not read file")
	}
	var c = Configuration{}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Println("Could not unmarshal yaml")
		panic("Could not unmarshal yaml")
	}

	return c, nil
}
