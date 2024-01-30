package pad

import (
	"testing"
)

func Test_parseHosts(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        int
	}

	var tests = []Test{
		{"Extra spaces and CRs", "h1   e.com:6501\r\nh2   e.com:6502\r\n\r\n h3   e.com:6503", 3},
		{"Extra spaces and CRs with Tabs", "h1\te.com:6501\t\r\nh2\t   e.com:6502\r\nh\r\n h3   \te.com:6503\r\n", 3},
		{"Empty element", "h1\te.com:6501\t\r\nh2\t\r\nh\r\n h3   \te.com:6503\r\n", 2},
		{"Empty string", "", 0},
	}

	for _, test := range tests {
		if got := parseHosts(test.input); len(got) != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
func Test_cleanStringSlice(t *testing.T) {

	type Test struct {
		description string
		input       []string
		want        int
	}

	var tests = []Test{
		//{"Extra spaces and CRs", "h1   e.com:6501\r\nh2   e.com:6502\r\n\r\n h3   e.com:6503\r\n", 3},
		//{"Extra spaces and CRs with Tabs", "h1\te.com:6501\t\r\nh2\t   e.com:6502\r\nh\r\n h3   \te.com:6503\r\n", 3},
		{"Three empty elements", []string{"", "", "", "1", "2"}, 2},
		{"One empty elements", []string{"1", "2", "", "3", "4"}, 4},
		{"All empty elements", []string{"", "", "", "", ""}, 0},
	}
	for _, test := range tests {
		if got := cleanStringSlice(test.input); len(got) != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
