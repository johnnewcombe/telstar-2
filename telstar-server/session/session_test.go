package session

import (
	"github.com/johnnewcombe/telstar-library/types"
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func Test_isempty(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_push(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_pop(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_pops(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_getCurrentFrameIdFromCache(t *testing.T) {

	sessionId := "1234567890"
	ClearCache(sessionId)
	CreateSession(sessionId, types.User{})
	PushHistory(sessionId, "101a")

	type Test struct {
		description string
		input       string
		want        string
	}

	var tests = []Test{
		{"ASCII 30h", "1234567890", "101a"},
	}

	for _, test := range tests {
		if got, ok := PeekHistory(sessionId); got != test.want || !ok || IsHistoryEmpty(sessionId) {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}
