package pad

import (
	"strings"
	"testing"
)

func Test_getProfileValue(t *testing.T) {

	type Test struct {
		description  string
		inputProfile Profile
		inputIndex   int
		wantValue    int
		wantOk       bool
	}

	var profile = Profile{"Profile T1",
		[]ProfileValue{
			{1, 0xff, "PAD Recall"},
			{2, 0x7f, "Echo"},
		},
	}

	var tests = []Test{
		{
			"Get index 1",
			profile,
			1,
			0xFF,
			true,
		},
		{
			"Get index 2",
			profile,
			2,
			0x7F,
			true,
		},
		{
			"Get index 3",
			profile,
			3,
			0,
			false,
		},
	}

	for _, test := range tests {
		if got, ok := getProfileValue(&test.inputProfile, test.inputIndex); got != test.wantValue || ok != test.wantOk {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_setProfileValue(t *testing.T) {

	type Test struct {
		description  string
		inputProfile Profile
		inputIndex   int
		inputValue   int
		wantValue    int
		wantOk       bool
	}

	var profile = Profile{"Profile T1",
		[]ProfileValue{
			{1, 0x7f, "PAD Recall"},
			{2, 0x7f, "Echo"},
		},
	}

	var tests = []Test{
		{
			"Set index 1",
			profile,
			1,
			0xFF,
			0xFF,
			true,
		},
		{
			"Set index 2",
			profile,
			2,
			0xFE,
			0xFE,
			true,
		},
		{
			"Set index 3 (Bad Index)",
			profile,
			3,
			0xFD,
			0,
			false,
		},
	}

	for _, test := range tests {
		if ok := setProfileValue(&test.inputProfile, test.inputIndex, test.inputValue); ok != test.wantOk {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
		// read the change if the previous test indicated a change
		if test.wantOk {
			if test.inputProfile.values[test.inputIndex-1].value != test.wantValue {
				t.Errorf(TEST_ERROR_MESSAGE, test.description)
			}
		}
	}
}

func Test_setProfile(t *testing.T) {

	type Test struct {
		description  string
		inputProfile string
		wantOk       bool
	}

	var tests = []Test{
		{"Profile P8 (ucase)", "P8", true},
		{"Profile 7 (lcase)", "7", false},
		{"Profile P6 (lcase)", "p6", false},
		{"Profile P (lcase)", "p", false},
		{"Profile P7 (lcase)", "p7", true},
		{"Profile P8 (lcase)", "p8", true},
	}
	for _, test := range tests {
		if ok := setProfile(test.inputProfile); ok != test.wantOk {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
		// read the change if the previous test indicated a change
		if test.wantOk {
			if strings.ToLower(test.inputProfile) != strings.ToLower(currentProfile.name) {
				t.Errorf(TEST_ERROR_MESSAGE, test.description)
			}
		}
	}
}
