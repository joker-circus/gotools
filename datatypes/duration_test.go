package datatypes

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	var cases = []struct {
		in  string
		out time.Duration

		expectedString string
	}{
		{
			in:             "0",
			out:            0,
			expectedString: "0s",
		}, {
			in:             "0w",
			out:            0,
			expectedString: "0s",
		}, {
			in:  "0s",
			out: 0,
		}, {
			in:  "324ms",
			out: 324 * time.Millisecond,
		}, {
			in:  "3s",
			out: 3 * time.Second,
		}, {
			in:  "-3s",
			out: -3 * time.Second,
		}, {
			in:  "5m",
			out: 5 * time.Minute,
		}, {
			in:  "1h",
			out: time.Hour,
		}, {
			in:  "4d",
			out: 4 * 24 * time.Hour,
		}, {
			in:  "4d1h",
			out: 4*24*time.Hour + time.Hour,
		}, {
			in:             "14d",
			out:            14 * 24 * time.Hour,
			expectedString: "2w",
		}, {
			in:  "3w",
			out: 3 * 7 * 24 * time.Hour,
		}, {
			in:             "3w2d1h",
			out:            3*7*24*time.Hour + 2*24*time.Hour + time.Hour,
			expectedString: "23d1h",
		}, {
			in:  "10y",
			out: 10 * 365 * 24 * time.Hour,
		},
	}

	for _, c := range cases {
		d, err := ParseDuration(c.in)
		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}
		if time.Duration(d) != c.out {
			t.Errorf("Expected %v but got %v", c.out, d)
		}
		expectedString := c.expectedString
		if c.expectedString == "" {
			expectedString = c.in
		}
		if d.String() != expectedString {
			t.Errorf("Expected duration string %q but got %q", c.in, d.String())
		}
	}
}

func TestDuration_UnmarshalText(t *testing.T) {
	var cases = []struct {
		in  string
		out time.Duration

		expectedString string
	}{
		{
			in:             "0",
			out:            0,
			expectedString: "0s",
		}, {
			in:             "0w",
			out:            0,
			expectedString: "0s",
		}, {
			in:  "0s",
			out: 0,
		}, {
			in:  "324ms",
			out: 324 * time.Millisecond,
		}, {
			in:  "3s",
			out: 3 * time.Second,
		}, {
			in:  "5m",
			out: 5 * time.Minute,
		}, {
			in:  "1h",
			out: time.Hour,
		}, {
			in:  "4d",
			out: 4 * 24 * time.Hour,
		}, {
			in:  "4d1h",
			out: 4*24*time.Hour + time.Hour,
		}, {
			in:             "14d",
			out:            14 * 24 * time.Hour,
			expectedString: "2w",
		}, {
			in:  "3w",
			out: 3 * 7 * 24 * time.Hour,
		}, {
			in:             "3w2d1h",
			out:            3*7*24*time.Hour + 2*24*time.Hour + time.Hour,
			expectedString: "23d1h",
		}, {
			in:  "10y",
			out: 10 * 365 * 24 * time.Hour,
		},
	}

	for _, c := range cases {
		var d Duration
		err := d.UnmarshalText([]byte(c.in))
		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}
		if time.Duration(d) != c.out {
			t.Errorf("Expected %v but got %v", c.out, d)
		}
		expectedString := c.expectedString
		if c.expectedString == "" {
			expectedString = c.in
		}
		text, _ := d.MarshalText() // MarshalText returns hardcoded nil
		if string(text) != expectedString {
			t.Errorf("Expected duration string %q but got %q", c.in, d.String())
		}
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	var cases = []struct {
		in  string
		out time.Duration

		expectedString string
	}{
		{
			in:             `"0"`,
			out:            0,
			expectedString: `"0s"`,
		}, {
			in:             `"0w"`,
			out:            0,
			expectedString: `"0s"`,
		}, {
			in:  `"0s"`,
			out: 0,
		}, {
			in:  `"324ms"`,
			out: 324 * time.Millisecond,
		}, {
			in:  `"3s"`,
			out: 3 * time.Second,
		}, {
			in:  `"5m"`,
			out: 5 * time.Minute,
		}, {
			in:  `"1h"`,
			out: time.Hour,
		}, {
			in:  `"4d"`,
			out: 4 * 24 * time.Hour,
		}, {
			in:  `"4d1h"`,
			out: 4*24*time.Hour + time.Hour,
		}, {
			in:             `"14d"`,
			out:            14 * 24 * time.Hour,
			expectedString: `"2w"`,
		}, {
			in:  `"3w"`,
			out: 3 * 7 * 24 * time.Hour,
		}, {
			in:             `"3w2d1h"`,
			out:            3*7*24*time.Hour + 2*24*time.Hour + time.Hour,
			expectedString: `"23d1h"`,
		}, {
			in:  `"10y"`,
			out: 10 * 365 * 24 * time.Hour,
		},
		{
			in:  `"289y"`,
			out: 289 * 365 * 24 * time.Hour,
		},
	}

	for _, c := range cases {
		var d Duration
		err := json.Unmarshal([]byte(c.in), &d)
		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}
		if time.Duration(d) != c.out {
			t.Errorf("Expected %v but got %v", c.out, d)
		}
		expectedString := c.expectedString
		if c.expectedString == "" {
			expectedString = c.in
		}
		bytes, err := json.Marshal(d)
		if err != nil {
			t.Errorf("Unexpected error on marshal of %v: %s", d, err)
		}
		if string(bytes) != expectedString {
			t.Errorf("Expected duration string %q but got %q", c.in, d.String())
		}
	}
}

func TestParseBadDuration(t *testing.T) {
	var cases = []string{
		"1",
		"1y1m1d",
		"-1w",
		"1.5d",
		"d",
		"294y",
		"200y10400w",
		"107675d",
		"2584200h",
		"",
	}

	for _, c := range cases {
		_, err := ParseDuration(c)
		if err == nil {
			t.Errorf("Expected error on input %s", c)
		}

	}
}

func BenchmarkParseDuration(b *testing.B) {
	const data = "30s"

	for i := 0; i < b.N; i++ {
		_, err := ParseDuration(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
