package session

import (
	"errors"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/types"
	"time"
)

// global session data, holds data for ALL connected users
type sessionData map[string]Session

var sessions = make(sessionData)

type Session struct {
	SessionId        string
	ConnectionNumber int    // number of the connection for current listener
	IPAddress        string // remote IPAddress
	History          []string
	FrameCache       map[string]types.Frame
	User             types.User
}

func CreateSession(sessionId string, user types.User, connectionNumber int, ipAddress string) Session {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "CreateSession")
	}

	// create the session and add to global session data
	s := Session{SessionId: sessionId}

	// frame cache
	s.FrameCache = make(map[string]types.Frame)

	// current user
	s.User = user

	// connection number of current listener run
	s.ConnectionNumber = connectionNumber

	// ip address of remote connection
	s.IPAddress = ipAddress

	// store session
	sessions[sessionId] = s

	// return session
	return s

}

func DeleteSession(sessionId string) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "DeleteSession")
	}

	delete(sessions, sessionId)
}

func UpdateCurrentUser(sessionId string, user types.User) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "UpdateCurrentUser")
	}

	s := sessions[sessionId]
	s.User = user
	sessions[sessionId] = s

}

func GetCurrentUser(sessionId string) types.User {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "GetCurrentUser")
	}

	return sessions[sessionId].User
}

func AddFrameToCache(sessionId string, frame types.Frame) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "AddFrameToCache")
	}

	s := sessions[sessionId]
	s.FrameCache[frame.GetPageId()] = frame
	sessions[sessionId] = s
}

func GetFrameFromCache(sessionId string, pageId string) (types.Frame, error) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "GetFrameFromCache")
	}

	frame := sessions[sessionId].FrameCache[pageId]
	if len(frame.PID.FrameId) == 0 {
		return frame, errors.New("frame not found")
	}
	return frame, nil
}

func ClearCache(sessionId string) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "ClearCache")
	}

	s := sessions[sessionId]
	s.FrameCache = make(map[string]types.Frame)
	sessions[sessionId] = s
}

// IsHistoryEmpty check if stack is empty
func IsHistoryEmpty(sessionId string) bool {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "IsHistoryEmpty")
	}

	s := sessions[sessionId]
	return len(s.History) == 0
}

// PushHistory a new value onto the stack
func PushHistory(sessionId string, pageId string) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "PushHistory")
	}

	// append the page and update the global data
	s := sessions[sessionId]
	s.History = append(s.History, pageId) // Simply append the new value to the end of the stack
	sessions[s.SessionId] = s
}

// PopHistoryN removes and return nth elements from the top of the stack. Return false if stack is empty after n Pops.
func PopHistoryN(sessionId string, n int) (string, bool) {
	var (
		pageId string
		ok     bool
	)

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "getFrame")
	}

	for p := 0; p < n; p++ {
		pageId, ok = PopHistory(sessionId)
	}
	return pageId, ok
}

// PopHistory removes and return top element of stack. Return false if stack is empty.
func PopHistory(sessionId string) (string, bool) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "PopHistory")
	}

	s := sessions[sessionId]
	if IsHistoryEmpty(sessionId) {
		return "", false
	} else {
		index := len(s.History) - 1     // Get the index of the top most element.
		element := (s.History)[index]   // Index into the slice and obtain the element.
		s.History = (s.History)[:index] // Remove it from the stack by slicing it off.

		// update the global variable
		sessions[s.SessionId] = s
		return element, true
	}
}

// PeekHistory return top element of stack, without removing it. Return false if stack is empty.
func PeekHistory(sessionId string) (string, bool) {

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "PeekHistory")
	}

	s := sessions[sessionId]

	if IsHistoryEmpty(sessionId) {
		return "", false
	} else {
		index := len(s.History) - 1     // Get the index of the top most element.
		return (s.History)[index], true // Index into the slice and obtain the element.
	}
}
