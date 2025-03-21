package netClient

import (
	"bufio"
	"context"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/utils"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/session"
	"github.com/johnnewcombe/telstar/synchronisation"
	"net"
	"time"

	"github.com/johnnewcombe/telstar-library/logger"
)

/*
The process is to start to create a common context with a cancel function before starting
two xfer go routines. Each is added to a common wait group.

This routine then waits indefinitely within the defer() function for both go routines to end

If either go routine sees a DLE or detects a network error then a cancel is issued
this will stop both go routines.
*/

func Connect(conn net.Conn, url string, connectionNumber int, baudRate int, initBytes []byte, settings config.Config) bool {

	// connect to remote host
	var (
		logPreAmble string
	)
	remoteConn, err := net.Dial("tcp", url)

	logPreAmble = utils.FormatLogPreAmble(session.GetSessionCount(), connectionNumber, "")

	if err != nil {
		logger.LogError.Printf("%sTCP connection to %s, failed. Error: %s", logPreAmble, url, err)
		return false
	}

	logPreAmble = utils.FormatLogPreAmble(session.GetSessionCount(), connectionNumber, utils.GetIpAddress(conn))
	logger.LogInfo.Printf("%sTCP connection made to %s.", logPreAmble, url)

	// signal channel, this is used to allow the xfer goroutine(s) to signal that a DLE
	// Data Link Escape has been detected (typically Ctrl P, see settings)
	//var dleSignal = make(chan bool)

	// wait group to wait for goroutines to complete before exiting this function
	waitGroup := synchronisation.WaitGroupWithCount{}

	// Context for cancellation, this is passed to the go routines which check the associate Done channel
	// to see if they should exit. The var 'cancel' is a function that can be called to set the Done channel
	// appropriately and tell the goroutine to stop. This is typically only done when a DLE has been
	// signalled by one of the goroutines
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {

		// no channels to close as only the don channel is used, this gets closed by cancel()
		// wait for transfer go routines to end
		waitGroup.Wait()

		logger.LogInfo.Printf("%sClosing TCP connection to %s.", logPreAmble, url)

		if err = remoteConn.Close(); err != nil {
			logger.LogError.Printf("%v", err)
		}

		if globals.Debug {
			// close the remote connection at the end
			logger.TimeTrack(time.Now(), "Connect")
		}
	}()

	///////////////////////////////////////////////////
	// Transfer Go Routine Definition
	///////////////////////////////////////////////////
	transfer := func(src net.Conn, dst net.Conn) {
		var (
			initDisabled   bool
			reader         *bufio.Reader
			timeoutCounter int
		)

		// the last thing we do is tell the wait group when we have completed
		defer func() {

			if err = src.SetReadDeadline(time.Time{}); err != nil {
				logger.LogError.Printf("%s%v", logPreAmble, err)
			}

			waitGroup.Done()
			logger.LogInfo.Printf("%sTransfer routine exiting [%s to %s].",
				logPreAmble,
				src.RemoteAddr().String(),
				dst.RemoteAddr().String())

			if globals.Debug {
				logger.TimeTrack(time.Now(), "xfer")
			}

		}()

		// create a new buffered reader
		reader = bufio.NewReader(src)

		// this function reads from the remote host and passes data to the PAD user
		// conn is the connection to the pad,
		// remoteConn is the connection to the gateway host
		logger.LogInfo.Printf("%sTransfer routine starting [%s to %s].",
			logPreAmble,
			src.RemoteAddr().String(),
			dst.RemoteAddr().String())

		for {

			// process any requested cancellation by checking the Done channel of the context
			// the Done channel is set in the parent function in response to this function
			// indicating a DLE (Data Link Escape) on the dleSignal channel. See below
			select {
			case <-ctx.Done():
				// ctx is telling us to stop, probably because we have sent it a DLE signal on the
				// dleSignal channel (see below).
				logger.LogInfo.Printf("%sTransfer cancelled [%s to %s].",
					logPreAmble,
					utils.GetIpAddress(src),
					utils.GetIpAddress(dst))

				// reset the read timeout that was added. The timeout meant that the read from conn would
				// only wait 100ms before spinning round to try again
				// this allowed the go routine to check the cancellation channel every 100ms or so.
				// now that the go routines have both finished, this can be put back to zero (blocking)
				//src.SetReadDeadline(time.Time{})
				return

			default:
			}

			// this means that the read from conn will only wait 500ms before spinning round to try again
			// this allows the go routine to check the cancellation channel every 500ms or so.
			// this will be put back once the gateway go routines have completed
			if err = src.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
				logger.LogError.Printf("%s%v", logPreAmble, err)

				// cancelling here may violate the channel closing rules but in practice this was far
				// more reliable than signalling to the parent via another channel to issue the cancel()
				cancel()

				// we just go round again to check the done channel
				continue
			}

			start := time.Now()
			ok, inputByte := readByte(reader)

			if !ok {
				// each !ok read should be 500ms since the last, anything shorter means that
				// the connection has been broken
				if time.Since(start).Milliseconds() < 50 {

					logger.LogInfo.Printf("%sClient has disconnected.", logPreAmble)

					//error, the connection has clearly been broken
					cancel()
					// we just go round again to check the done channel
					continue
				}

				// this will run every 500ms if there is no activity
				// could time-ou connection after 5 mins
				timeoutCounter++
				if timeoutCounter > settings.Server.GatewayTimeout*2 { // TODO Add gateway timeout counter

					logger.LogInfo.Printf("%sInactivity timeout exceeded.", logPreAmble)

					// no activity on the  src connection so report a DLE
					// so that this goroutine and the goroutine handling the other direction
					// gets closed
					//dleSignal <- true
					cancel()

					// we just go round again to check the done channel
					continue
				}

				// check for EOF
				if inputByte == 0xFF {

					logger.LogError.Printf("%sEOF detected, simulating DLE.", logPreAmble)

					// this means the src connection is closed so report a DLE
					// so that this goroutine and the goroutine handling the other direction
					// gets closed
					//dleSignal <- true
					cancel()
					continue
				}
				// Input byte error detected, ignoring.
				time.Sleep(100 * time.Millisecond)
				continue

			}

			// char received so reset timeout counter
			timeoutCounter = 0

			if globals.Debug {
				logger.LogInfo.Printf("%sByte '%02x' received from %s.", logPreAmble, inputByte, src.RemoteAddr().String())
			}

			write := func(data []byte) {
				if _, err := dst.Write(data); err != nil {
					logger.LogError.Printf("%sWRITE: %s", logPreAmble, err)

					cancel()

				} else if globals.Debug {
					logger.LogInfo.Printf("%sByte Sent to %s.", dst.RemoteAddr().String())
				}
			}

			if inputByte == settings.Server.DLE {

				cancel()

			} else {
				//
				if len(initBytes) > 0 && !initDisabled {
					// if we have init bytes e.g. from the minitel parser, send them
					write(initBytes)
					// we only do this once
					initDisabled = true
				}
				write([]byte{inputByte})

			}

			// baud rate simulation
			time.Sleep(time.Duration(baudRate))
		}

	}

	// transfer from remote to connected user and visa versa
	// The go routine will signal to the wait group when it has been completed
	waitGroup.Add(2)
	go transfer(conn, remoteConn)
	go transfer(remoteConn, conn)

	// return OK after waiting for all go routines to end see defer() function
	return true
}

// TODO this same function is defined in PAD, SERVER and Net Client
func readByte(reader *bufio.Reader) (bool, byte) {

	// get a byte
	inputByte, err := reader.ReadByte()
	if err != nil {
		return false, inputByte
	}
	return true, inputByte
}
