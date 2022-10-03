package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Database   DatabaseConfig                `yaml:"database"`
	WebService WebServiceConfig              `yaml:"webservice"`
	Parsers    []map[string]ParserYamlConfig `yaml:"parsers"`
}
type YamlParser struct {
	Name string
	Url  string `yaml:"url"`
}
type DatabaseConfig struct {
	Url string `yaml:"url"`
}

type WebServiceConfig struct {
	Port uint16 `yaml:"port"`
}

type ParserYamlConfig struct {
	Url string `yaml:"url"`
}

type ParserConfig struct {
	Name string
	Url  string
}

type Parsers []ParserConfig

type Message struct {
	Environments map[string]models `yaml:"Environments"`
}
type models map[string][]Model
type Model struct {
	AppType     string `yaml:"app-type"`
	ServiceType string `yaml:"service-type"`
}

func main() {

	yfile, err := os.ReadFile("aggregator.yaml")

	if err != nil {

		log.Fatal(err)
	}

	var yamlConfig YamlConfig

	err2 := yaml.Unmarshal(yfile, &yamlConfig)

	if err2 != nil {

		log.Fatal(err2)
	}

	parsers := make(Parsers, len(yamlConfig.Parsers))

	for i, v := range yamlConfig.Parsers {
		for name, parser := range v {
			parsers[i] = ParserConfig{Name: name, Url: parser.Url}
		}
	}

	fmt.Println(yamlConfig)
	fmt.Println(parsers)
}
