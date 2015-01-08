package main

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Sizes map[Size]int

type Store map[string]string

type Handler struct {
	Sizes   []Sizes
	Formats []Format
	Store   []Store
}

// Config contains the configuration options for the service
type Config map[string]Handler

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
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
