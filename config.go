package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type SizePair struct {
	Size  Size
	Value int
}

type Store map[string]string

type Handler struct {
	Sizes   []SizePair
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

func (sp *SizePair) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmap map[string]int
	if err := unmarshal(&tmap); err != nil {
		return err
	}
	for k, v := range tmap {
		switch Size(k) {
		case Square, Max, MaxHeight, MaxWidth:
			sp.Size = Size(k)
			sp.Value = v
		default:
			return errors.New(fmt.Sprintf("Invalid size in config: %s", k))
		}
	}
	return nil
}

// func (f *Format) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	var fm string
// 	if err := unmarshal(&fm); err != nil {
// 		return err
// 	}
// 	format := Format(fm)
// 	switch format {
// 	case JPG, PNG, WebP:
// 		print(format)
// 		f = &format
// 	default:
// 		return errors.New(fmt.Sprintf("Invalid file format in config: %s", fm))
// 	}
// 	return nil
// }
