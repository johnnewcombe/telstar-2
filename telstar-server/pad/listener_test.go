package pad

import (
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func Test_isAplphNumeric(t *testing.T) {

	type Test struct {
		description string
		input       byte
		want        bool
	}

	var tests = []Test{
		{"ASCII 30h", 0x30, true},
		{"ASCII 31h", 0x31, true},
		{"ASCII 32h", 0x32, true},
		{"ASCII 33h", 0x33, true},
		{"ASCII 34h", 0x34, true},
		{"ASCII 35h", 0x35, true},
		{"ASCII 36h", 0x36, true},
		{"ASCII 37h", 0x37, true},
		{"ASCII 38h", 0x38, true},
		{"ASCII 39h", 0x39, true},
		{"ASCII 40h", 0x40, true},
		{"ASCII 29h", 0x29, true},
		{"ASCII 07h", 0x07, false},
		{"ASCII 08h", 0x08, true},
		{"ASCII 09h", 0x09, false},
		{"ASCII 1fh", 0x1f, false},
		{"ASCII 20h", 0x20, true},
		{"ASCII 7fh", 0x7f, true},
		{"ASCII 80h", 0x80, false},
	}
	for _, test := range tests {
		if got := isAlphaNumeric(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
