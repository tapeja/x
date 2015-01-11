package x

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
    - jpg
    - png
  store:
    - local: /mnt/gallery`,
			expected: &Config{
				"android-tablets": Handler{
					[]Size{
						Size{MaxWidth, 500},
						Size{Max, 1000},
					},
					[]Format{JPG, PNG},
					[]StoreConfig{
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
			t.Errorf("Expected error %q, got %q for %q", tc.err, err, tc.yaml)
		}
	}
}
