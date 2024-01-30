package pad

import "testing"

func Test_isStringInSlice(t *testing.T) {

	type Test struct {
		description string
		input       string
		slice       []string
		want        bool
	}

	list := []string{"help", "profile", "profiles", "exit"}

	var tests = []Test{
		{"'help'", "help", list, true},
		{"'profile'", "profile", list, true},
		{"' help '", " help ", list, false},
		{"'dog'", "dog", list, false},
		{"Space", " ", list, false},
		{"Empty input", "", list, false},
		{"Empty slice", "", []string{""}, false},
	}
	for _, test := range tests {
		if got := isStringInSlice(test.input, list); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_parseCommand(t *testing.T) {

	type Test struct {
		description string
		input       string
		wantIsValid bool
	}

	var tests = []Test{
		//{"'help'", "help", true},
		//{"'help me'", "help me", true},
		{"' help '", " help ", true},
		{"'dog'", "dog", false},
		{"Space", " ", false},
		{"Empty input", "", false},
	}
	for _, test := range tests {
		if got := parseCommand(test.input); got.isValid != test.wantIsValid {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
func Test_isProfileSetting(t *testing.T) {
	type Test struct {
		description string
		input       string
		want        bool
	}
	var tests = []Test{
		{"PROFILE = P0", "profile = p0", true},
		{" PRO FILE = P 0", " pro file = p 0", true},
		{"PROFILE=P0", "profile=p0", true},
		{"PROFILE=0", "profile=0", false},
		{"PROFILE=", "profile=", false},
		{"PROFIL=P0", "profil=P0", false},
		{"PROFILE1=P0", "profile1=p0", false},
	}

	for _, test := range tests {
		if got := isProfileSetting(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_isParamSetting(t *testing.T) {
	type Test struct {
		description string
		input       string
		want        bool
	}
	var tests = []Test{
		{"P2=3", "P2=3", true},
		{" P 2 = 3", "p 2 = 3", true},
		{"P2=3", "p2=3", true},
		{"P2=0", "p2=0", true},
		{"P2=", "p2=", false},
		{"P0=2", "p0=2", false},
		{"P2=334", "p2=334", true},
		{"P22=334", "p22=334", true},
		{"P223=334", "p222=334", false},
		{"P22=3334", "p222=334", false},
	}

	for _, test := range tests {
		if got := isParamSetting(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
func Test_parseProfile(t *testing.T) {
	type Test struct {
		description string
		input       string
		want        string
	}
	var tests = []Test{
		{"PROFILE=P3", "profile=p3", "p3"},
		{"PROFILE=3", "profile=3", ""},
		{"PROFILE=P33", "profile=p33", ""},
		{"PROFIL=P3", "profil=p3", ""},
	}

	for _, test := range tests {
		if got := parseProfile(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
func Test_parseParam(t *testing.T) {
	type Test struct {
		description string
		input       string
		wantValue1  string
		wantValue2  string
	}
	var tests = []Test{
		{"P2=3", "p2=3", "2", "3"},
		{"P2=0", "p2=0", "2", "0"},
		{"P2=", "p2=", "", ""},
		{"P2=334", "p2=334", "2", "334"},
		{"P0=2", "p0=2", "", ""},
	}

	for _, test := range tests {
		if got1, got2 := parseParam(test.input); got1 != test.wantValue1 || got2 != test.wantValue2 {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
