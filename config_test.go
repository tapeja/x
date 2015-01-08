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
	}{
		{
			yaml: `android-tablets:
  sizes:
    - max_width: 500
    - max: 1000
  formats:
    - jpeg
  store:
    - local: /mnt/gallery`,
			expected: &Config{
				"android-tablets": Handler{
					[]Sizes{
						Sizes{Size(MaxWidth): 500},
						Sizes{Size(Max): 1000},
					},
					[]Format{Format("jpeg")},
					[]Store{
						{
							"local": "/mnt/gallery",
						},
					},
				},
			},
			err: nil,
		},
	}
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
