package main

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

type Handler struct {
	Sizes   map[Size]int
	Formats []Format
	Stores  []url.URL
}

type yamlHandler struct {
	Sizes   map[string]int
	Formats []string
	Stores  []string
}

type yamlConfig map[string]yamlHandler

// Config contains the configuration options for the service
type Config struct {
	Handlers map[string]Handler
}

// NewConfig creates new configuration
func NewConfig() *Config {
	return &Config{}
}

// Load reads the configuration file from the provided file name
func (c *Config) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.LoadStream(f)
}

// LoadStream reads the configuration from the provided reader
func (c *Config) LoadStream(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	var yc yamlConfig
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &yc); err != nil {
		return err
	}
	return nil
}
