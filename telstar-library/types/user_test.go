package types

import (
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func Test_IsInScope(t *testing.T) {

	type Test struct {
		description       string
		inputUserBasePage int
		inputPageNumber   int
		output            bool
	}

	var tests = []Test{
		{"Page 0, BasePage 0", 0, 0, true},
		{"Page 999, BasePage 0", 0, 999, true},
		{"Page 0, BasePage Empty", 0, 0, false},
		{"Page 0, BasePage 100", 100, 0, false},
		{"Page 100, BasePage 100", 100, 100, true},
		{"Page 200, BasePage 200", 200, 200, true},
		{"Page 200, BasePage 2013", 2013, 200, false},
		{"Page 201, BasePage 2013", 2013, 201, false},
		{"Page 2013, BasePage 2013", 2013, 2013, true},
		{"Page 20131, BasePage 2013", 2013, 20131, true},
		{"Page 201312345, BasePage 2013", 2013, 201312345, true},
		{"Page 100, BasePage 999999999", 999999999, 100, false},
	}

	for _, test := range tests {

		user := User{"7777777777", "", "", test.inputUserBasePage, true, true, false, true}

		if got := user.IsInScope(test.inputPageNumber); got != test.output {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
func Test_IsGuest(t *testing.T) {

	type Test struct {
		description       string
		inputUserBasePage int
		output            bool
	}

	var tests = []Test{
		{"Test 1", 100, false},
		{"Test 2", 0, true},
		{"Test 3", 0, true},
		{"Test 4", -1, false},
	}

	for _, test := range tests {

		user := User{"7777777777", "", "", test.inputUserBasePage, true, true, false, true}

		if got := user.IsGuest(); got != test.output {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

