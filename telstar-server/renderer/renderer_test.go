package renderer

import (
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar/config"
	"context"
	"testing"
	"time"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func Test_render(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_getGreeting(t *testing.T) {

	type Test struct {
		description    string
		inputTime      string
		outputGreeting string
	}

	var tests = []Test{
		{"", "10:00", "GOOD MORNING"},
		{"", "00:00", "GOOD MORNING"},
		{"", "13:00", "GOOD AFTERNOON"},
		{"", "12:00", "GOOD AFTERNOON"},
		{"", "17:00", "GOOD EVENING"},
		{"", "17:01", "GOOD EVENING"},
	}

	for _, test := range tests {

		tm, _ := time.Parse("15:04", test.inputTime)
		//fmt.Printf("%v - %v", t,err)

		if got := getGreeting(tm); got != test.outputGreeting {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_renderHeader(t *testing.T) {

	t.Error("Test not implemented!")
	sessionId:="1234567890"

	conn := MockConn{} // mock connection object
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type Test struct {
		description     string
		inputFrameType  string
		inputHideCost   bool
		inputHidePageId bool
		inputCls        bool
		inputParity     bool
		inputCursor     bool
		inputHeader     string
		inputPageNo     int
		inputFrameId    string
		//output          []byte
		outputLen int
	}

	// table to test the buffer
	var tests = []Test{
		{"Default Header", "information", false, false, true, false, false,
			"", 0, "a",
			//[]byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			40},

		{"Specified Header", "information", false, false, true, false, false,
			"[Y]MICRONETn 800 (C)", 800, "a",
			//[]byte{12, 27, 67, 77, 73, 67, 82, 79, 78, 69, 84, 110, 32, 56, 48, 48, 32, 40, 67, 41, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			40},

		// note that if frametype is initial and cls is true, metadata will be returned based on current date etc
		{"Meta Header", "initial", false, false, true, false, false,
			"", 0, "a",
			//[]byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			54},

		{"Meta Header no CLS", "information", false, false, false, false, false,
			"", 0, "a",
			//[]byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			40},

		{"Default Header hidden page Id", "information", false, true, true, false, false,
			"", 0, "a",
			//[]byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			40},

		{"Default Header hidden cost", "information", true, false, true, false, false,
			"", 0, "a",
			//[]byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			40},

		{"Default Header hidden cost and pageId", "information", true, true, true, false, false,
			"", 0, "a",
			//[]byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32},
			40},
	}

	for _, test := range tests {

		conn.buffer = []byte{}

		frame := getEmptyFrame()
		frame.PID.PageNumber = test.inputPageNo
		frame.PID.FrameId = test.inputFrameId
		frame.HeaderText = test.inputHeader
		frame.Cursor = test.inputCursor
		frame.FrameType = test.inputFrameType

		options := getTestRenderOptions()
		//options.ClearScreen = test.inputCls
		//options.baudRate
		//options.hasFollowOnFrame

		settings := getTestConfig()
		settings.Server.HideCost = test.inputHideCost
		settings.Server.HidePageId = test.inputHidePageId

		renderHeader(ctx, &conn, &frame, sessionId, settings, options)

		//if getDisplayLenTest(string(conn.buffer)) != test.outputLen {
		//	t.Errorf("%s: Length of header is incorrect", test.description)
		//}
		//if !compareByteSlices(conn.buffer, test.output) {
		//	t.Errorf("%s: Header data does not match that expected\r\nwant:%v\r\n got:%v", test.description, test.output, conn.buffer)
		//}
	}

	//h1 := []byte{12, 27, 66, 84, 27, 65, 69, 27, 70, 76, 27, 68, 83, 27, 71, 84, 27, 69, 65, 27, 67, 82, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 20, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32}
	//h2 := []byte{12, 27, 67, 77, 73, 67, 82, 79, 78, 69, 84, 110, 32, 56, 48, 48, 32, 40, 67, 41, 32, 32, 32, 32, 32, 32, 32, 20, 56, 48, 48, 97, 32, 32, 32, 32, 32, 32, 32, 32, 48, 112, 32}
}
func Test_renderContent(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_RenderSystemMessage(t *testing.T) {

	t.Error("Test not implemented!")
	/*
	   	conn := MockConn{} // mock connection object
	   	ctx, cancel := context.WithCancel(context.Background())
	   	defer cancel()

	   	type Test struct {
	   		description   string
	   		inputSettings config.Config
	   		inputOptions  RenderOptions
	   		inputFrame    dal.Frame
	   		output        []byte
	   	}

	   	// table to test the buffer
	   	var tests = []Test{
	   		{"Page not found", getTestConfig(), getTestRenderOptions(), getEmptyFrame(),
	   			[]byte("\x1e\x0b\x1b\x44\x1b\x5d\x1b\x43Page not Found")},
	   	}

	   	for _, test := range tests {
	   		RenderSystemMessage(ctx, &conn, &test.inputFrame, test.inputSettings, dal.User{}, test.inputOptions)//; !utils.AreEqualByteSlices(conn.buffer, test.output) ||
	   //			err != nil {
	   //			t.Errorf(TEST_ERROR_MESSAGE, test.description)
	   //		}
	   	}
	*/

}

func Test_renderBuffer(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_padTextLeft(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_padTextRight(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_getLen(t *testing.T) {
	t.Error("Test not implemented!")
}

func getTestConfig() config.Config {
	settings := config.Config{}
	settings.Server.Strings.DefaultPageNotFoundMessage = ""
	settings.Server.Strings.DefaultNavMessage= "[B][n][Y]Select item or[W]*page# : "
	settings.Server.Strings.DefaultHeaderText = "[G]T[R]E[C]L[B]S[W]T[M]A[Y]R"
	return settings
}
func getTestRenderOptions() RenderOptions {
	options := RenderOptions{}
	//options.ClearScreen = true
	options.BaudRate = globals.BAUD_RATE
	options.HasFollowOnFrame = false
	return options
}

func getEmptyFrame() types.Frame {
	return types.Frame{}
}

