package main

import "testing"

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func Test_createDefaultRoutingTable(t *testing.T) {

	type Test struct {
		description string
		input       int
		inputRoute  int
		wantRoute   int
	}

	var tests = []Test{
		{"Hash Route", 8001221, 10, 800},
		{"0 Route", 8001221, 0, 80012210},
		{"1 Route", 8001221, 1, 80012211},
		{"2 Route", 8001221, 2, 80012212},
		{"3 Route", 8001221, 3, 80012213},
		{"4 Route", 8001221, 4, 80012214},
		{"5 Route", 8001221, 5, 80012215},
		{"6 Route", 8001221, 6, 80012216},
		{"7 Route", 8001221, 7, 80012217},
		{"8 Route", 8001221, 8, 80012218},
		{"9 Route", 8001221, 9, 80012219},
	}

	// run tests
	for _, test := range tests {
		if got := createDefaultRoutingTable(test.input); got[test.inputRoute] != test.wantRoute {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_checkArgCount(t *testing.T) {

	type Test struct {
		description string
		inputArgs   []string
		inputCount  int
		wantOk      bool
	}

	var tests = []Test{
		{"", []string{"hello", "world", "dog"}, 2, false},
		{"", []string{"hello", "world", "priMarY"}, 2, true},
		{"", []string{"hello", "world", "secONDARY"}, 2, true},
		{"", []string{"hello", "world", "foo", "bar"}, 5, false},
		{"", []string{"hello", "world", "foo", "bar"}, 3, false},
	}
	// run tests
	for _, test := range tests {
		if got := checkArgCount(test.inputArgs, test.inputCount); got != test.wantOk {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_getPidFromFileName(t *testing.T) {

	type Test struct {
		description string
		input       string
		wantPageNo  int
		wantFrameId string
		wantOk      bool
	}

	var tests = []Test{
		{"filename somedirectory/1234a.edit.tf", "tmp/1234a.edit.tf", 1234, "a", true},
		{"filename 1234a.editt.f", "1234a.editt.f", 0, "", false},
		{"filename a.edit.tf", "a.edit.tf", 0, "", false},
		{"filename 1234a.edittf", "1234a.edittf", 0, "", false},
		{"filename 1234a.edittf", "1234a.edittf", 0, "", false},
		{"filename 1234a.edit.tf", "1234a.edit.tf", 1234, "a", true},
	}

	for _, test := range tests {
		if got, ok := getPidFromFileName(test.input); got.PageNumber != test.wantPageNo ||
			got.FrameId != test.wantFrameId ||
			ok != test.wantOk {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
