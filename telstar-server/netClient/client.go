package netClient

import (
	"bufio"
	"context"
	"net"
	"sync"
	"time"

	"github.com/johnnewcombe/telstar-library/logger"
)

func Connect(conn net.Conn, url string, dataLinkEscape byte, baudRate int, initBytes []byte) bool {

	// connect to remote host
	remoteConn, err := net.Dial("tcp", url)
	if err != nil {
		logger.LogError.Printf("TCP connection to %s, failed. Error: %s", url, err)
		return false
	}
	logger.LogInfo.Printf("TCP connection made to %s.", url)

	// close the remote connection at the end
	defer remoteConn.Close()

	// signal channel, this is used to allow the xfer goroutine(s) to signal that a DLE
	// Data Link Escape has been detected (typically Ctrl P, see settings)
	var dleSignal = make(chan bool)

	// Context for cancellation, this is passed to the go routines which check the associate Done channel
	// to see if they should exit. The var 'cancel' is a function that can be called to set the Done channel
	// appropriately and tell the goroutine to stop. This is typically only done when a DLE has been
	// signalled by one of the goroutines
	ctx, cancel := context.WithCancel(context.Background())

	// wait group to wait for goroutines to complete before exiting this function
	var waitgroup sync.WaitGroup

	// xfer from remote to pad user, note that we add 1 to the wait group and pass a pointer to the waitgroup
	// to the goroutine. The gorouting will signal to the waitgroupo when it has been completed
	waitgroup.Add(1)
	go xfer(ctx, &waitgroup, conn, remoteConn, dataLinkEscape, baudRate, initBytes, dleSignal)

	// xfer from pad user to remote, note that we add 1 to the wait group and pass a pointer to the waitgroup
	//	// to the goroutine. The gorouting will signal to the waitgroupo when it has been completed
	waitgroup.Add(1)
	go xfer(ctx, &waitgroup, remoteConn, conn, dataLinkEscape, baudRate, initBytes, dleSignal)

	// this will block until one of the xfer functions sends a DLE detected message
	select {
	case dle := <-dleSignal:
		if dle {
			// we have the dle  from one of the xfer routines so tell them both to stop
			logger.LogInfo.Println("DLE Received, cancelling gateway activity.")
			cancel()
		}
	}

	// tel the waitgroup to wait for all goroutines to complete
	waitgroup.Wait()

	// return OK
	return true

}

func xfer(ctx context.Context, waitgroup *sync.WaitGroup, src net.Conn, dst net.Conn, dataLinkEscape byte, baudRate int, initBytes []byte, dleSignal chan bool) {

	var (
		initDisabled bool
		reader       *bufio.Reader
	)

	// the last thing we do is tell the waitgroup when we have completed
	defer waitgroup.Done()
	defer src.SetReadDeadline(time.Time{})

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
			logger.LogInfo.Printf("Transfer cancelled [%s to %s].",
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

		// this means that the read from conn will only wait 100ms before spinning round to try again
		// this allows the go routine to check the cancellation channel every 100ms or so.
		// this will be put back once the gateway go routines have completed
		src.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		ok, inputByte := readByte(reader)

		if !ok {
			// this means the src connection is closed so report a DLE
			// so that this goroutine and the goroutine handling the other direction
			// gets closed

			// check for EOF
			if inputByte == 0xFF {
				logger.LogError.Print("EOF detected, simulating DLE.")
				dleSignal <- true
			}
			// Input byte error detected, ignoring.
			time.Sleep(100 * time.Millisecond)
			continue

		}
		logger.LogInfo.Printf("Byte '%02x' received from %s.", inputByte, src.RemoteAddr().String())

		write := func(data []byte) {
			if _, err := dst.Write(data); err != nil {
				logger.LogError.Printf("WRITE: %s", err)
			} else {
				logger.LogInfo.Printf("Byte Sent to %s.", dst.RemoteAddr().String())
			}
		}

		if inputByte == dataLinkEscape {

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
			/*
				if _, err := dst.Write([]byte{inputByte}); err != nil {
					logger.LogError.Printf("WRITE: %s", err)
				} else {
					logger.LogInfo.Printf("Byte Sent to %s.", dst.RemoteAddr().String())
				}
			*/
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
	logger.LogInfo.Println("character read:", inputByte)
	return true, inputByte
}

/*
func readByte(conn net.Conn) (bool, byte) {

	result := true
	inputByte := make([]byte, 1, 1)

	count, err := conn.Read(inputByte)

	if err != nil {
		if err != io.EOF {
			logger.LogWarn.Println("read error:", err)
			return false, 0
		} else {
			return false, 255 // EOF
		}
		result = false
	}
	//logger.LogInfo.Println("ReadByte Done")

	if count > 0 {
		return result, inputByte[0]
	} else {
		return false, 0
	}
}

*/
