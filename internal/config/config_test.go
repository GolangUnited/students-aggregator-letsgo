package config

import (
	"testing"
)

const (
	configFilePath = "../../tests/configs/config/config.yaml"
)

func TestRead(t *testing.T) {

	c := NewConfig()

	err := c.SetDataFromFile(configFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if c.Database.Name != "stub" {
		t.Errorf("incorrect database url \"%v\", expected \"%s\"", c.Database.Name, "stub")
	}

	if c.Database.Url != "stub://localhost:22222/" {
		t.Errorf("incorrect database url \"%v\", expected \"%s\"", c.Database.Url, "stub://localhost:22222/")
	}

	if c.WebService.Port != 8080 {
		t.Errorf("incorrect webservice port %d, expected %d", c.WebService.Port, 8080)
	}

	if len(c.Parsers) != 3 {
		t.Errorf("incorrect parsers count %d, expected %d", len(c.Parsers), 3)
	}

	for _, p := range c.Parsers {
		switch p.Name {
		case "github":
			if p.URL != "https://github.com/golang/go/tags" {
				t.Errorf("incorrect parser \"%s\" url \"%s\", expected \"%s\"", p.Name, p.URL, "https://github.com/golang/go/tags")
			}
		case "go.dev":
			if p.URL != "https://go.dev/blog" {
				t.Errorf("incorrect parser \"%s\" url \"%s\", expected \"%s\"", p.Name, p.URL, "https://go.dev/blog")
			}
		case "medium.com":
			if p.URL != "https://medium.com/_/graphql" {
				t.Errorf("incorrect parser \"%s\" url \"%s\", expected \"%s\"", p.Name, p.URL, "https://medium.com/_/graphql")
			}
		default:
			t.Errorf("unknown parser \"%s\"", p.Name)
		}
	}
}
