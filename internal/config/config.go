package config

import (
	"fmt"
	"os"

	aggregator "github.com/indikator/aggregator_lets_go/internal/aggregator/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	webservice "github.com/indikator/aggregator_lets_go/internal/webservice/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	data       []byte
	Aggregator aggregator.Config
	Database   db.Config
	WebService map[string]webservice.Config
	Parsers    []parser.Config
}

func NewConfig() *Config {
	return &Config{}
}

type parserYamlConfig struct {
	URL      string            `yaml:"url"`
	LogLevel logLevel.LogLevel `yaml:"logLevel"`
}

type webserviceYamlConfig struct {
	Handle   string            `yaml:"handle"`
	Port     uint16            `yaml:"port"`
	LogLevel logLevel.LogLevel `yaml:"logLevel"`
}

type yamlConfig struct {
	Aggregator aggregator.Config             `yaml:"aggregator"`
	Database   db.Config                     `yaml:"database"`
	Parsers    []map[string]parserYamlConfig `yaml:"parsers"`
	Handles  []map[string]webserviceYamlConfig `yaml:"webservice"`
}

func (c *Config) SetDataFromFile(fileName string) error {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	return c.SetData(data)
}

func (c *Config) SetData(data []byte) error {
	c.data = data
	return nil
}

func (c *Config) Read() (err error) {
	if len(c.data) == 0 {
		return fmt.Errorf("data to read not found")
	}

	var yc yamlConfig

	err = yaml.Unmarshal(c.data, &yc)

	if err != nil {
		return
	}

	c.Aggregator = yc.Aggregator
	c.Database = yc.Database
	c.WebService = make(map[string]webservice.Config)
	c.Parsers = make([]parser.Config, len(yc.Parsers))

	for _, v := range yc.Handles {
		for name, h := range v {
			c.WebService[name] = webservice.Config{Handle: h.Handle, Port: h.Port, LogLevel: h.LogLevel}
		}
	}

	for i, v := range yc.Parsers {
		for name, p := range v {
			c.Parsers[i] = parser.Config{Name: name, URL: p.URL}
		}
	}

	return nil
}
