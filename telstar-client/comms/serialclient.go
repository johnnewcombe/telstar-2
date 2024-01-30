package comms

import (
	"context"
	"fmt"
	"go.bug.st/serial.v1"
	"io"
	"log"
	"sync"
	"time"
)

// see https://pkg.go.dev/go.bug.st/serial.v1#Mode

type SerialClient struct{}

var (
	port serial.Port
)

func (c *SerialClient) Open(endpoint Endpoint) error {

	var (
		parity   serial.Parity
		dataBits int
		err      error
		mode     serial.Mode
	)

	if err = c.Close(); err != nil {
		return err
	}

	if endpoint.Serial.Parity {
		parity = serial.EvenParity
		dataBits = 7
	} else {
		parity = serial.NoParity
		dataBits = 8
	}

	mode = serial.Mode{
		BaudRate: endpoint.Serial.Baud,
		Parity:   parity,
		DataBits: dataBits,
		StopBits: serial.OneStopBit,
	}

	// Note that he configuration can be changed at any time with the SetMode function:
	if port, err = serial.Open(endpoint.Serial.Port, &mode); err != nil {
		port = nil
		return err
	}

	//port.SetDTR(true)
	//port.SetRTS(true)
	m, err := port.GetModemStatusBits()
	print(m)

	time.Sleep(500 * time.Millisecond)

	// Modem Init
	if len(endpoint.Serial.ModemInit) > 0 {
		if err = c.Write([]byte(fmt.Sprintf("%s\r\n", endpoint.Serial.ModemInit))); err != nil {
			return err
		}
	}

	return nil
}

func (c *SerialClient) Close() error {

	if port != nil {
		if err := port.Close(); err != nil {
			return err
		}
		port = nil
	}
	return nil
}

func (c *SerialClient) Write(byt []byte) error {

	var err error

	if port != nil {
		if _, err = port.Write(byt); err != nil {
			return err
		}
		//print(string(byt))
	}
	return nil

}

func (c *SerialClient) Read(ctx context.Context, wg *sync.WaitGroup, f funcDef) {

	defer wg.Done()

	for {

		// process any requested cancellation by checking the Done channel of the context
		// note that the resdByte function blocks and so when in that state cancellation
		// wont happen unless there is data coming in or the port is closed
		select {
		case <-ctx.Done():
			// ctx is telling us to stop
			log.Println("SerialClient.Read() goroutine cancelled")
			return

		default:
		}

		if port != nil {
			ok, inputByte := c.readByte()

			// DEBUG
			fmt.Printf("OK:%v, byte:%d\r\n", ok, inputByte)

			f(ok, inputByte)
		} else {
			// no need to rush as the connection isn't open
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func (c *SerialClient) readByte() (bool, byte) {

	result := true
	inputByte := make([]byte, 1, 1)

	// get a byte
	count, err := port.Read(inputByte)

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

	if len(inputByte) > 0 {
		return result, inputByte[0]
	} else {
		return false, 0
	}
}
