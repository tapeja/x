package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	var testCases = []struct {
		yaml     string
		expected *Config
		err      error
	}{}
	for _, tc := range testCases {
		c := NewConfig()
		err := c.LoadStream(strings.NewReader(tc.yaml))
		if !reflect.DeepEqual(tc.expected, c) {
			t.Errorf("Expected %q, got %q for %q", tc.expected, c, tc.yaml)
		}
		if err != tc.err {
			t.Errorf("Expected erro %q, got %q for %q", tc.err, err, tc.yaml)
		}
	}
}
