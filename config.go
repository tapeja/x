package main

import (
	"os",
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {

}

func NewConfig() {
	return &Config{}
}

func (c *Config) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	c.LoadStream(f)
}

func (c *Config) LoadStream(r io.Reader) error {
	data, err := ReadAll(r)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}
