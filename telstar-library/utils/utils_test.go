package utils

import (
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func Test_IsValidPageId(t *testing.T) {

	// define tests
	type Test struct {
		description       string
		input             string
		wantIsValidPageId bool
	}

	var tests = []Test{
		{"PageId 123a", "123a", true},
		{"PageId 999z", "999z", true},
		{"PageId 999999999a", "999999999a", true},
		{"PageId a", "a", false},
		{"PageId empty", "", false},
		{"PageId 1", "1", false},
		{"PageId 999", "999", false},
		{"PageId 999!", "999!", false},
		{"PageId 9999999999a", "9999999999a", false},
	}

	// run tests
	for _, test := range tests {
		if got := IsValidPageId(test.input); got != test.wantIsValidPageId {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_IsValidFrameId(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IsValidRoutingTable(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IsValidPageNumber(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IsValidUserId(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_CreateDefaultRoutingTable(t *testing.T) {

	// define tests
	type Test struct {
		description     string
		inputPageNumber int
		wantTable       []int
	}
	var tests = []Test{
		{"Test 1", 221, []int{2210, 2211, 2212, 2213, 2214, 2215, 2216, 2217, 2218, 2219, 22}},
		{"Test 1", 1, []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 0}},
	}
	// run tests
	for _, test := range tests {
		if got := CreateDefaultRoutingTable(test.inputPageNumber); !compareRoutingTable(got, test.wantTable) {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_ConvertPidToPageId(t *testing.T) {

	// define tests
	type Test struct {
		description     string
		inputPageNumber int
		inputFrameId    string
		wantPageId      string
	}

	var tests = []Test{
		{"PID 123/a", 123, "a", "123a"},
		{"PID 999/z", 999, "z", "999z"},
		{"PID 999999999/a", 999999999, "a", "999999999a"},
	}

	// run tests
	for _, test := range tests {
		if got, err := ConvertPidToPageId(test.inputPageNumber, test.inputFrameId); got != test.wantPageId ||
			err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_ConvertPageIdToPID(t *testing.T) {

	// define tests
	type Test struct {
		description    string
		inputPageId    string
		wantPageNumber int
		wantFrameId    string
	}

	var tests = []Test{
		{"PageId 123a", "123a", 123, "a"},
		{"PageId 999z", "999z", 999, "z"},
		{"PageId 999999999a", "999999999a", 999999999, "a"},
	}

	// run tests
	for _, test := range tests {
		if gotPageNumber, gpotFrameId, err := ConvertPageIdToPID(test.inputPageId); gotPageNumber != test.wantPageNumber ||
			gpotFrameId != test.wantFrameId ||
			err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_IsNumeric(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IsAlphaNumeric(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IsGraphicCharacter(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}
func Test_IsGraphicColor(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IsAlphaColour(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_GetFollowOnFrameId(t *testing.T) {

	type Test struct {
		description string
		input       rune
		want        rune
	}
	var tests = []Test{
		{"FrameId a", 'a', 'b'},
		{"FrameId c", 'c', 'd'},
		{"FrameId z", 'z', 'a'},
	}
	for _, test := range tests {
		if got, err := GetFollowOnFrameId(test.input); err != nil || got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_GetFollowOnPageId(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        string
	}
	var tests = []Test{
		{"PageId 99a", "99a", "99b"},
		{"PageId 121c", "121c", "121d"},
		{"PageId 121z", "121z", "1210a"},
		{"PageId 999999999y", "999999999y", "999999999z"},
	}
	for _, test := range tests {
		if got, err := GetFollowOnPageId(test.input); err != nil || got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_GetFollowOnPID(t *testing.T) {

	type Test struct {
		description     string
		inputPageNumber int
		inputFrameId    rune
		wantPageNumber  int
		wantFrameId     rune
	}
	var tests = []Test{
		{"PageId 99a", 99, 'a', 99, 'b'},
		{"PageId 121c", 121, 'c', 121, 'd'},
		{"PageId 121z", 121, 'z', 1210, 'a'},
		{"PageId 999999999y", 999999999, 'y', 999999999, 'z'},
	}
	for _, test := range tests {
		if got1, got2, err := GetFollowOnPID(test.inputPageNumber, test.inputFrameId); err != nil ||
			got1 != test.wantPageNumber || got2 != test.wantFrameId {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_ParseDataType(t *testing.T) {

	type Test struct {
		description string
		input       string
		outputType  string
		outputParam int
	}

	var tests = []Test{
		{"Test 1", "edit.tf,5", "edit.tf", 5},
		{"Test 2", "edittf,Hello", "edittf", 0},
		{"Test 3", "markup,-5", "markup", 0},
		{"Test 4", "vraw,55", "vraw", 0},
		{"Test 5", "traw,23", "traw", 0},
		{"Test 6", "edit.tf,22", "edit.tf", 22},
	}

	for _, test := range tests {

		if gotType, gotParam := ParseDataType(test.input); gotType != test.outputType || gotParam != test.outputParam {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_IntToByte(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_IntToBytes(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_CheckPasswordStrength(t *testing.T) {

	type Test struct {
		description string
		input       string
		pass        bool
	}

	var tests = []Test{
		{"Test 1", "1234", false},
		{"Test 2", "January16", false},
		{"Test 2", "January16-6$", true},
		{"Test 3", "1234567890", false},
		{"Test 4", "Hello World", false},
		{"Test 5", "01e1954c-1595-45bf-bb6f-fa3f6a5e5ec4", true},
	}

	for _, test := range tests {

		if err := CheckPasswordStrength(test.input); test.pass == (err != nil) {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_GetDateTimeFromUnixData(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_GetDateFromUnixData(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_GetTimeFromUnixData(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func compareRoutingTable(tableA []int, tableB []int) bool {

	if tableA == nil || tableB == nil || len(tableA) != len(tableB) || len(tableA) != 11 {
		return false
	}

	for i := 0; i < 11; i++ {
		if tableA[i] != tableB[i] {
			return false
		}
	}
	return true
}
