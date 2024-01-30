package pad

import (
	"context"
	"net"
	"time"

	"bitbucket.org/johnnewcombe/telstar-library/logger"
)

const(
	MAX_INACTIVE_CRS = 3
)

var promptCount int // keeps track of the number of CR's sent without a command being executed

// this is needed to ensure that we don't scroll the screen.
func render(ctx context.Context, conn net.Conn, data string, baudRate int) {

	displayPrompt := func() {
		if promptCount >= MAX_INACTIVE_CRS {
			writeString(conn, WELCOME_MESSAGE, baudRate)
			writeString(conn, PROMPT, baudRate)
			promptCount = 0
		} else {
			writeString(conn, PROMPT, baudRate)
			promptCount++
		}
		commandBuffer = ""

	}

	for _, b := range []byte(data) {

		// if we get a CLS/HOME then we can reset the number of allowed CRs entered by the user
		if b == 0xc {
			promptCount = 0
		}
		// process any requested cancellation
		select {
		case <-ctx.Done():
			// channel has a true, so cancel
			displayPrompt()
			return // channel closed so cancel
		default:
		}
		// write data
		if _, err := conn.Write([]byte{b}); err != nil {
			logger.LogError.Print(err)
		}
		// baud rate simulation
		time.Sleep(time.Duration(baudRate))
	}
	displayPrompt()
	return
}

func writeString(conn net.Conn, data string, baudRate int) {

	for _, b := range []byte(data) {

		// if we get a CLS/HOME then we can reset the number of allowed CRs entered by the user
		if b == 0xc {
			promptCount = 0
		}

		// write data
		if _, err := conn.Write([]byte{b}); err != nil {
			logger.LogError.Print(err)
		}
		// baud rate simulation
		time.Sleep(time.Duration(baudRate))
	}

}
