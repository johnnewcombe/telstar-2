package routing

import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar/session"
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
	SESSIONID = "7f1208f9-b3ab-4844-be1d-3e0d6bcdeeea"
)

func Test_RoutingRequest_isValid(t *testing.T) {

	type Test struct {
		description string
		input       RouterRequest
		wantIsValid bool
	}

	var tests = []Test{
		{"valid RouterRequest", RouterRequest{32, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, ""}, true},
		{"short routing table", RouterRequest{32, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},SESSIONID},false},                  // wrong number of elements in routing tbl
		{"very Short routing table", RouterRequest{32, "0a", false, []int{0},SESSIONID}, false},                                        // wrong number of elements in routing tbl
		{"empty routing table", RouterRequest{32, "0a", false, []int{},SESSIONID}, false},                                              // no routing table
		{"bad pageId", RouterRequest{32, "9999999999a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID}, false},              // bad pageId
		{"bad input byte", RouterRequest{0, "9a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID}, false},                    // bad inputByte
		{"bad routing table entry", RouterRequest{32, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9999999999, 10},SESSIONID}, false}, // bad routing entry
		{"bad routing table entry", RouterRequest{32, "0a", false, []int{0, 1111111111, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID}, false}, // bad routing entry

	}

	// run tests
	for _, test := range tests {
		if got := test.input.isValid(); got != test.wantIsValid {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_processRouting(t *testing.T) {

	type Test struct {
		description       string
		input             RouterRequest
		output            RouterResponse
		wantStatus        int
		wantRoutingBuffer string
		wantNewPageId     string
		wantHistoryPage   bool
		wantImmediateMode bool
	}

	var tests = []Test{

		{
			"immidiate mode with input byte = 'A'",
			RouterRequest{0x41, "0a", false, []int{90, 91, 92, 93, 94, 95, 6, 97, 98, 99, 0},SESSIONID},
			RouterResponse{0, "", "", false, true},
			InvalidCharacter, "", "", false, true,
		},

		{
			"immidiate mode with input byte = '1'",
			RouterRequest{0x31, "0a", false, []int{90, 91, 92, 93, 94, 95, 6, 97, 98, 99, 0},SESSIONID},
			RouterResponse{0, "", "", false, true},
			ValidPageRequest, "", "91a", false, true,
		},

		{
			"immidiate mode with input byte = '#'",
			RouterRequest{0x5f, "12a", false, []int{90, 91, 92, 93, 94, 95, 6, 97, 98, 99, 0},SESSIONID},
			RouterResponse{0, "", "", false, true},
			ValidPageRequest, "", "12b", false, true,
		},
	}

	// run tests
	for _, test := range tests {

		if ProcessRouting(&test.input, &test.output); //test.output.currentPageId != test.wantCurrentPageId ||
		//test.output.followOnPageId != test.wantFollowOnPageId ||
		test.output.HistoryPage != test.wantHistoryPage ||
			test.output.ImmediateMode != test.wantImmediateMode ||
			test.output.NewPageId != test.wantNewPageId ||
			test.output.RoutingBuffer != test.wantRoutingBuffer ||
			test.output.Status != test.wantStatus {

			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_RoutingResponse_trimBuffer(t *testing.T) {

	type Test struct {
		description string
		input       RouterResponse
		wantBuffer  string
	}

	var tests = []Test{

		{"buffer empty", RouterResponse{1, "", "", false, true}, ""},
		{"buffer contains '*_'", RouterResponse{2, "*_", "", false, true}, ""},
		{"buffer contains '*'", RouterResponse{3, "*", "", false, true}, ""},
		{"buffer contains '***'", RouterResponse{4, "***", "", false, true}, ""},
		{"buffer contains '***_'", RouterResponse{5, "***_", "", false, true}, ""},
		{"buffer contains '***___'", RouterResponse{6, "***___", "", false, true}, ""},
		{"buffer contains '_'", RouterResponse{7, "_", "", false, true}, ""},
		{"buffer contains '___'", RouterResponse{8, "___", "", false, true}, ""},
		{"buffer contains '*_*_'", RouterResponse{9, "*_*_", "", false, true}, "_*"},
	}

	// run tests
	for _, test := range tests {
		if test.input.trimBuffer(); test.input.RoutingBuffer != test.wantBuffer {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_truncateData(t *testing.T) {

	type Test struct {
		description string
		input       RouterRequest
		output      RouterResponse
		wantBuffer  string
	}

	// table to test the buffer
	var tests = []Test{
		{
			"empty buffer",
			RouterRequest{0x41, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"",
		},
		{
			"3 chars, non-immediate = no trim",
			RouterRequest{0x44, "1a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "ABC", "", false, false},
			"ABC",
		},
		{
			"3 chars, immediate = no trim",
			RouterRequest{0x34, "2a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, true},
			"123",
		},
		{
			"11 chars, non-immediate = no trim",
			RouterRequest{0x41, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123456789AB", "", false, false},
			"123456789AB",
		},
		{
			"11 chars, immediate = no trim",
			RouterRequest{0x41, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123456789AB", "", false, true},
			"123456789AB",
		},
		{
			"12 chars, non-immediate = trimmed all but initial '*'",
			RouterRequest{0x42, "4a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "1234567890ABC", "", false, false},
			"",
		},
		{
			"12 chars, immediate = trimmed",
			RouterRequest{0x42, "5a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123456789ABC", "", false, true},
			"",
		},
	}

	// run tests
	for _, test := range tests {
		if truncateData(&test.input, &test.output); test.output.RoutingBuffer != test.wantBuffer {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_preProcessInputByte(t *testing.T) {

	// tests the correct selection of immediate mode
	type Test struct {
		description       string
		input             RouterRequest
		output            RouterResponse
		wantImmediateMode bool
		wantBuffer        string
	}

	// table to test the buffer
	var tests = []Test{
		{
			"asterisk when in immediate mode should set immediate mode to false (buffer mode)",
			RouterRequest{0x2a, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			false, "",
		},
		{

			"asterisk when in buffer mode should keep immediate mode false (buffer mode)",
			RouterRequest{0x2a, "1a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, false},
			false, "", // buffer should stay empty
		},
		{
			"asterisk when buffer not empty and in buffer mode should keep",
			RouterRequest{0x2a, "2a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, false},
			false, "", //immediate mode false (buffer mode) but clear the buffer
		},
		{
			"non-asterisk when in immediate mode should not set immediate mode to false",
			//non-asterisk when in immediate mode should not set immediate mode to false
			//the buffer should not be updated, only backspace causes a buffer change
			RouterRequest{0x41, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, true},
			true, "123", //the buffer should not be updated, only backspace causes a buffer change
		},
		{
			"non-asterisk when not in immediate mode should not set immediate mode to true",
			RouterRequest{0x34, "4a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, false},
			false, "123", //the buffer should not be updated, only backspace causes a buffer change
		},

		{
			"backspace with no buffer immediate mode",
			RouterRequest{0x08, "5a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			true, "", // no buffer change
		},
		{
			"last char removed from the buffer",
			RouterRequest{0x08, "6a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123A", "", false, true},
			true, "123", // backspace with buffer immediate mode
		},
		{
			"backspace with no buffer non-immediate mode",
			RouterRequest{0x08, "7a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, false},
			false, "",
		},
		{
			"backspace with buffer non-immediate mode",
			RouterRequest{0x08, "8a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123A", "", false, false},
			false, "123",
		},
		{
			"normal characters",
			RouterRequest{0x41, "9a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "3a", false, true},
			true, "123",
		},
	}

	// run tests
	for _, test := range tests {
		if preProcessInputByte(&test.input, &test.output); test.output.ImmediateMode != test.wantImmediateMode ||
			test.output.RoutingBuffer != test.wantBuffer {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_selectReloadFrame(t *testing.T) {

	type Test struct {
		description       string
		input             RouterRequest
		output            RouterResponse
		wantNewPageId     string
		wantImmediateMode bool
	}

	// table to test the buffer
	var tests = []Test{
		{
			"frame 0a, immediate mode",
			RouterRequest{0, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"0a", true, // set immediate mode to true and set new page to 0a
		},
		{
			"frame 1a, immediate mode",
			RouterRequest{0, "1a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, false},
			"1a", true, // set immediate mode to true and set new page to 1a
		},
		{
			"bad pageId, immediate mode",
			RouterRequest{0, "", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"", true,
		},
	}

	// run tests
	for _, test := range tests {
		if selectReloadFrame(&test.input, &test.output); test.output.NewPageId != test.wantNewPageId ||
			test.output.ImmediateMode != test.wantImmediateMode {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_selectPreviousFrame(t *testing.T) {

	type Test struct {
		description       string
		input             RouterRequest
		output            RouterResponse
		wantBuffer        string
		wantImmediateMode bool
		wantHistoryPage   bool
		wantHistoryEmpty  bool
	}

	// create a user and a session
	session.CreateSession(SESSIONID, types.User{"0","somehash","John",0,true, true,true,true})
	session.PushHistory(SESSIONID, "7a")
	session.PushHistory(SESSIONID, "8a")
	session.PushHistory(SESSIONID, "9a")

	// table to test the buffer
	var tests = []Test{

		{
			"frame 0a, immediate mode, with history",
			RouterRequest{0, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"9a", true, true, false, // set immediate mode to true and set new page to 0a
		},
		{
			"frame 1a, immediate mode, with history",
			RouterRequest{0, "1a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, false},
			"8a", true, true, false, // set immediate mode to true and set new page to 1a
		},
		{
			"bad pageId, immediate mode, with history",
			RouterRequest{0, "", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"7a", true, true, true, // set immediate mode to true, set new pageId be invalid also
		},

		// History should be empty now

		{
			"frame 0a, immediate mode, no history",
			RouterRequest{0, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"3a", true, false, true, // set immediate mode to true and set new page to 0a
		},
		{
			"frame 1a, immediate mode, no history",
			RouterRequest{0, "4a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, false},
			"4a", true, false, true, // set immediate mode to true and set new page to 1a
		},
		{"bad pageId, immediate mode, no history",
			RouterRequest{0, "", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"", true, false, true, // set immediate mode to true, set new pageId be invalid also
		},
	}

	// run tests
	for _, test := range tests {
		if selectPreviousFrame(&test.input, &test.output); test.output.NewPageId != test.wantBuffer ||
			test.output.ImmediateMode != test.wantImmediateMode ||
			test.output.HistoryPage != test.wantHistoryPage ||
			session.IsHistoryEmpty(SESSIONID) != test.wantHistoryEmpty {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_selectFollowOnFrame(t *testing.T) {

	type Test struct {
		description         string
		input               RouterRequest
		output              RouterResponse
		wantFollowOnFrameId string
		wantStatus          int
	}
	var tests = []Test{
		{
			"invalid page request 0a",
			RouterRequest{0, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"0b", ValidPageRequest,
		},
		{
			"valid page request 0z",
			RouterRequest{0, "0z", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"00a", ValidPageRequest,
		},
		{
			"valid page request 3y",
			RouterRequest{0, "3y", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"3z", ValidPageRequest,
		},
		{
			"valid page request 5z",
			RouterRequest{0, "5z", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"50a", ValidPageRequest,
		},
		{
			"valid page request 999999999z",
			RouterRequest{0, "999999999z", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"", InvalidPageRequest, // this page is invalid
		},
	}

	// run tests
	for _, test := range tests {
		if selectFollowOnFrame(&test.input, &test.output); test.output.NewPageId != test.wantFollowOnFrameId ||
			test.output.Status != test.wantStatus {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_processImmediateMode(t *testing.T) {

	type Test struct {
		description   string
		input         RouterRequest
		output        RouterResponse
		wantNewPageId string
		wantStatus    int
	}

	// table to test the buffer
	var tests = []Test{

		{
			"invalid character 'A'",
			RouterRequest{0x41, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"", InvalidCharacter, //only numbers allowed in immediate mode
		},

		/* Calling processImmediateMode() when in Non-Immediate mode is invalid.
		{"",
			RouterRequest{0x31, "0a", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			RouterResponse{0, "ABC", "", "", false, false},
			"", VALID_PAGE_REQUEST,
		},
		{"",
			RouterRequest{0x32, "1a", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			RouterResponse{0, "ABC", "", "", false, false},
			"", VALID_PAGE_REQUEST,
		},
		*/

		{
			"valid page request '4'",
			// element 4 of the routing table will be returned
			RouterRequest{0x34, "2a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, true},
			"4a", ValidPageRequest,
		},

		{
			"valid page request '1'",
			// element 1 of the routing table will be returned
			RouterRequest{0x31, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, true},
			"1a", ValidPageRequest,
		},
		{
			"valid page request '9'",
			RouterRequest{0x39, "4a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123456789AB", "", false, true},
			"9a", ValidPageRequest, // element 1 of the routing table will be returned
		},
	}

	// run tests
	//FIXME Check newPageId also
	for _, test := range tests {
		if processImmediateMode(&test.input, &test.output); test.output.NewPageId != test.wantNewPageId ||
			test.output.Status != test.wantStatus {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_processBufferMode(t *testing.T) {
	type Test struct {
		description   string
		input         RouterRequest
		output        RouterResponse
		wantNewPageId string
		wantStatus    int
	}

	// table to test the buffer
	var tests = []Test{

		{
			"",
			// 'A' is ivalid as letters and numbers allowed in immediate mode
			RouterRequest{0x41, "0a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "", "", false, false},
			"", RouteMessageUpdated,
		},

		/* Calling processNonImmediateMode() when in Immediate mode is invalid.
		{"",
			RouterRequest{0x31, "0a", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			RouterResponse{0, "ABC", "", false, true},
			"", VALID_PAGE_REQUEST,
		},
		{"",
			RouterRequest{0x32, "1a", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			RouterResponse{0, "ABC", "", false, true},
			"", VALID_PAGE_REQUEST,
		},
		*/

		{
			"routing message updated '4'",
			RouterRequest{0x34, "2a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, false},
			"", RouteMessageUpdated, //  should just update the buffer, no new pageId yet
		},
		{
			"valid page request '#' with '123' in buffer",
			RouterRequest{0x5f, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "123", "", false, false},
			"123a", ValidPageRequest,
		},
		{
			"valid page request '#' with '*123' in buffer",
			RouterRequest{0x5f, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "*123", "123", false, false},
			"123a", ValidPageRequest,
		},
		{
			"invaid page request '#' with '*TEST' in buffer",
			RouterRequest{0x5f, "3a", false, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},SESSIONID},
			RouterResponse{0, "*TEST", "123", false, false},
			"TEST", InvalidPageRequest,
		},
	}

	// run tests
	for _, test := range tests {
		if processBufferMode(&test.input, &test.output); test.output.NewPageId != test.wantNewPageId ||
			test.output.Status != test.wantStatus {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_isNumeric(t *testing.T) {

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
		{"ASCII 40h", 0x40, false},
		{"ASCII 29h", 0x29, false},
		{"ASCII 08h", 0x08, false},
	}
	for _, test := range tests {
		if got := utils.IsNumeric(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

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
		if got := utils.IsAlphaNumeric(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_getFollowOnFrameId(t *testing.T) {

	type Test struct {
		description         string
		inputPageId         string
		wantFollowOnFrameId string
	}

	var tests = []Test{
		{"PageId 0a", "0a", "0b"},
		{"PageId 0z", "0z", "00a"}, // is this a valid page
		{"PageId 3a", "3a", "3b"},
		{"PageId 3z", "3z", "30a"},
		{"PageId 999999999z", "999999999z", "9999999990a"}, // this page is invalid though
	}

	// run tests
	for _, test := range tests {
		if got, err := GetFollowOnPageId(test.inputPageId); got != test.wantFollowOnFrameId || err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
