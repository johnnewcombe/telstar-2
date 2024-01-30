package comms

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// NetClient A struct  type that implements the CommunicationClient interface
type NetClient struct {}


var (
	conn net.Conn
)

func (c *NetClient) Open(endpoint Endpoint) error {

	var err error

	if err = c.Close(); err != nil {
		return err
	}
	// connect to remote host
	if conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d",
		endpoint.Address.Host,
		endpoint.Address.Port)); err != nil {

		conn = nil
		return err
	}
	if endpoint.Init.Telnet {
		if err = c.Write([]byte{255, 253, 3}); err != nil {
			return err
		}
	}
	if len(endpoint.Init.InitChars) > 0 {
		if err = c.Write(endpoint.Init.InitChars); err != nil {
			return err
		}
	}

	return nil
}

func (c *NetClient) Close() error {
	if conn != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		conn = nil
	}
	return nil
}

func (c *NetClient) Write(byt []byte) error {

	var err error

	if conn != nil {
		if _, err = conn.Write(byt); err != nil {
			return err
		}
	}
	return nil
}

func (c *NetClient) Read(ctx context.Context, wg *sync.WaitGroup, f funcDef) {

	defer wg.Done()

	for {

		// process any requested cancellation by checking the Done channel of the context
		select {
		case <-ctx.Done():
			// ctx is telling us to stop
			//log.Println("NetClient.Read() goroutine cancelled")

			return

		default:
		}

		if conn != nil {
			ok, inputByte := c.readByte()
			f(ok, inputByte)
		} else {
			// no need to rush as the connection isn't open
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func (c *NetClient) readByte() (bool, byte) {

	result := true
	inputByte := make([]byte, 1, 1)

	// get a byte
	if conn != nil {
		count, err := conn.Read(inputByte)

		if err != nil {
			if err != io.EOF {
				return false, 0
			} else {
				if err = c.Close(); err != nil {
					return false, 0
				}
			}
			result = false
		}
		if count == 0 {
			print(count)
		}
	}

	if len(inputByte) > 0 {
		return result, inputByte[0]
	} else {
		return false, 0
	}
}
