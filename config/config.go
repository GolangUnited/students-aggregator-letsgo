package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Url string `yaml:"url"`
}

type WebServiceConfig struct {
	Port uint16 `yaml:"port"`
}

type ParserConfig struct {
	Name string
	Url  string
}

type Config struct {
	Database   DatabaseConfig
	WebService WebServiceConfig
	Parsers    []ParserConfig
}

func NewConfig() *Config {
	return &Config{}
}

type parserYamlConfig struct {
	Url string `yaml:"url"`
}

type yamlConfig struct {
	Database   DatabaseConfig                `yaml:"database"`
	WebService WebServiceConfig              `yaml:"webservice"`
	Parsers    []map[string]parserYamlConfig `yaml:"parsers"`
}

func (c *Config) Read(fileName string) (err error) {
	file, err := os.ReadFile(fileName)

	if err != nil {

		return
	}

	var yc yamlConfig

	err = yaml.Unmarshal(file, &yc)

	if err != nil {

		return
	}

	c.Database = yc.Database
	c.WebService = yc.WebService
	c.Parsers = make([]ParserConfig, len(yc.Parsers))

	for i, v := range yc.Parsers {
		for name, parser := range v {
			c.Parsers[i] = ParserConfig{Name: name, Url: parser.Url}
		}
	}

	return nil
}
