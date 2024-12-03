package netClient

import (
	"bufio"
	"context"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/synchronisation"
	"net"
	"time"

	"github.com/johnnewcombe/telstar-library/logger"
)

/*
The process is to start to create a common context with a cancel function before starting
two xfer go routines. Each is added to a common wait group. The context is passed to each of
these routines. The wait group is also passed to each transfer routine so that the function
can indicate when it has completed.

This routine then waits indefinitely for a DLE on the dleSignal channel which is shared by
each xfer go routine. If either go routine sees a DLE it is signalled on the channel back
to this routine so that cancel can be set on the context. This should stop all go routines
*/

func Connect(conn net.Conn, url string, connectionNumber int, baudRate int, initBytes []byte, settings config.Config) bool {

	// connect to remote host
	var (
		userIp string
	)
	remoteConn, err := net.Dial("tcp", url)
	if err != nil {
		logger.LogError.Printf("%d:%s: TCP connection to %s, failed. Error: %s", connectionNumber, userIp, url, err)
		return false
	}

	// get remote IP Address
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		userIp = addr.IP.String()
	}
	logger.LogInfo.Printf("%d:%s: TCP connection made to %s.", connectionNumber, userIp, url)

	// signal channel, this is used to allow the xfer goroutine(s) to signal that a DLE
	// Data Link Escape has been detected (typically Ctrl P, see settings)
	var dleSignal = make(chan bool)

	// wait group to wait for goroutines to complete before exiting this function
	waitGroup := synchronisation.WaitGroupWithCount{}

	// Context for cancellation, this is passed to the go routines which check the associate Done channel
	// to see if they should exit. The var 'cancel' is a function that can be called to set the Done channel
	// appropriately and tell the goroutine to stop. This is typically only done when a DLE has been
	// signalled by one of the goroutines
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {

		cancel()
		waitGroup.Wait()

		logger.LogInfo.Printf("%d:%s: Closing TCP connection to %s.", connectionNumber, userIp, url)

		if err = remoteConn.Close(); err != nil {
			logger.LogError.Printf("%v", err)
		}

		if globals.Debug {
			// close the remote connection at the end
			logger.TimeTrack(time.Now(), "Connect")
		}
	}()

	// transfer from remote to connected user, note that we add 1 to the wait group and pass a pointer to the wait group
	// to the goroutine. The go routine will signal to the wait group when it has been completed
	waitGroup.Add(1)
	go xfer(ctx, &waitGroup, conn, remoteConn, connectionNumber, baudRate, initBytes, settings, dleSignal)

	// transfer from connected user to remote, note that we add 1 to the wait group and pass a pointer to the wait group
	// to the goroutine. The goroutine will signal to the wait group when it has been completed
	waitGroup.Add(1)
	go xfer(ctx, &waitGroup, remoteConn, conn, connectionNumber, baudRate, initBytes, settings, dleSignal)

	// this will block until one of the xfer functions sends a DLE detected message
	select {
	case dle := <-dleSignal:
		if dle {
			// we have the dle  from one of the xfer routines so tell them both to stop
			logger.LogInfo.Printf("%d:%s: DLE Received, cancelling gateway activity.", connectionNumber, userIp)
			// this routine will complete at the return below and the deferred function
			// will issue cancel and wait functions
		}
	}

	// return OK
	return true

}

func xfer(ctx context.Context, waitgroup *synchronisation.WaitGroupWithCount, src net.Conn, dst net.Conn, connectionNumber int, baudRate int, initBytes []byte, settings config.Config, dleSignal chan bool) {

	var (
		initDisabled   bool
		reader         *bufio.Reader
		err            error
		timeoutCounter int
		sourceIp       string
	)

	// the last thing we do is tell the wait group when we have completed
	defer func() {

		if err = src.SetReadDeadline(time.Time{}); err != nil {
			logger.LogError.Printf("%d:%s: %v", connectionNumber, sourceIp, err)
		}

		print()
		waitgroup.Done()

		if globals.Debug {
			logger.TimeTrack(time.Now(), "xfer")
		}

	}()

	// get remote IP Address
	if addr, ok := src.RemoteAddr().(*net.TCPAddr); ok {
		sourceIp = addr.IP.String()
	}

	// create a new buffered reader
	reader = bufio.NewReader(src)

	// this function reads from the remote host and passes data to the PAD user
	// conn is the connection to the pad,
	// remoteConn is the connection to the gateway host

	for {

		// process any requested cancellation by checking the Done channel of the context
		// the Done channel is set in the parent function in response to this function
		// indicating a DLE (Data Link Escape) on the dleSignal channel. See below
		select {
		case <-ctx.Done():
			// ctx is telling us to stop, probably because we have sent it a DLE signal on the
			// dleSignal channel (see below).
			logger.LogInfo.Printf("%d:%s: Transfer cancelled [%s to %s].",
				connectionNumber, sourceIp,
				src.RemoteAddr().String(),
				dst.RemoteAddr().String())

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
			logger.LogError.Printf("%d:%s: %v", connectionNumber, sourceIp, err)
			dleSignal <- true
		}

		ok, inputByte := readByte(reader)

		if !ok {

			// this will run every 500ms if there is no activity
			// could time-ou connection after 5 mins
			timeoutCounter++
			if timeoutCounter > settings.Server.GatewayTimeout*2 { // TODO Add gateway timeout counter

				logger.LogInfo.Printf("%d:%s: Inactivity timeout exceeded.", connectionNumber, sourceIp)

				// no activity on the  src connection so report a DLE
				// so that this goroutine and the goroutine handling the other direction
				// gets closed
				dleSignal <- true
			}

			// check for EOF
			if inputByte == 0xFF {

				logger.LogError.Printf("%d:%s: EOF detected, simulating DLE.", connectionNumber, sourceIp)

				// this means the src connection is closed so report a DLE
				// so that this goroutine and the goroutine handling the other direction
				// gets closed
				dleSignal <- true
			}
			// Input byte error detected, ignoring.
			time.Sleep(100 * time.Millisecond)
			continue

		}

		// char received so reset timeout counter
		timeoutCounter = 0

		if globals.Debug {
			logger.LogInfo.Printf("%d:%s: Byte '%02x' received from %s.", connectionNumber, sourceIp, inputByte, src.RemoteAddr().String())
		}

		write := func(data []byte) {
			if _, err := dst.Write(data); err != nil {
				logger.LogError.Printf("%d:%s: WRITE: %s", connectionNumber, sourceIp, err)

				// indicate on the dleSignal channel to the parent to send a cancel to this goroutine
				// stopping it here directly would violate the channel closing rules
				// signal DLE to sender so that sender can close using context.WithCancel()
				dleSignal <- true

			} else if globals.Debug {
				logger.LogInfo.Printf("%d:%s: Byte Sent to %s.", connectionNumber, sourceIp, dst.RemoteAddr().String())
			}
		}

		if inputByte == settings.Server.DLE {

			// indicate on the dleSignal channel to the parent to send a cancel to this goroutine
			// stopping it here directly would violate the channel closing rules
			// signal DLE to sender so that sender can close using context.WithCancel()
			dleSignal <- true

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

// TODO this same function is defined in PAD, SERVER and Net Client
func readByte(reader *bufio.Reader) (bool, byte) {

	// get a byte
	inputByte, err := reader.ReadByte()
	if err != nil {
		return false, inputByte
	}
	return true, inputByte
}
