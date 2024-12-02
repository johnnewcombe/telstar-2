package pad

import (
	"bufio"
	"context"
	"fmt"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/netClient"
	"github.com/johnnewcombe/telstar/server"
	"net"
	"strconv"
	"time"
)

const (
	WELCOME_MESSAGE = "\x0c\r\n\r\nTELSTAR VPAD\r\nFor assistance type 'HELP'.\r\n"
	PROMPT          = "\r\nPAD:00> "
	ERROR           = "\r\nERROR"
	HUH             = "\r\nHUH?"
	OK              = "\r\nOK"
	CUSTOM          = "Custom Profile"
	BAUD_2400       = 4167000 // nanoseconds

)

func Start(port int, settings config.Config, hosts map[string]string) error {

	logger.LogInfo.Printf("Starting PAD Server on port %d", port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Pad.Host, port))

	if err != nil {
		return err
	}

	for {
		// blocks until an incoming connection is made
		// when a connectioj is made it returns a net.Conn object
		conn, err := listener.Accept()
		if err != nil {
			logger.LogError.Print(err)
			continue
		}
		logger.LogInfo.Print("Incoming connection!")

		// handles one connection at a time
		go handleConn(conn, settings, hosts)

	}
}

func handleConn(conn net.Conn, settings config.Config, hosts map[string]string) {

	// handles one connection at a time
	var (
		reader          *bufio.Reader
		minitelResponse string
		minitelParser   server.MinitelParser
		err             error
		initBytes       []byte
	)

	defer closeConn(conn)

	// set current profile
	// TODO sort out case and constant
	if setProfile("p7") {
		logger.LogInfo.Printf("Setting current profile to %s.", DEFAULT_PROFILE)
	} else {
		logger.LogError.Printf("Unable to set current profile to %s.", DEFAULT_PROFILE)
	}

	// Telnet DO LINE MODE for testing only
	//writeBytes(conn, []byte{0xFF, 0xFD, 0x22, 0xFF, 0xFB, 0x01})
	baudRate := BAUD_2400

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writeString(conn, WELCOME_MESSAGE, baudRate)
	writeString(conn, PROMPT, baudRate)

	// create a new buffered reader
	reader = bufio.NewReader(conn)

	//loop as each character is received
	for {

		// get a byte
		ok, inputByte := readByte(reader)

		if !ok {
			//ignore the byte and spin around for the next one
			logger.LogError.Print("Input byte error detected, closing connection.")
			time.Sleep(100 * time.Millisecond)
			return
		}
		logger.LogError.Printf("Byte %s received.", string(inputByte))
		// cancel any existing rendering
		cancel()

		// pass through the Minitel parser, this will absorb any negotiation and
		// set minitelParser.minitelConnection to true if a Minitel negotiation was detected
		inputByte, minitelResponse = minitelParser.ParseMinitelDc(inputByte)

		// Minitel parser may need to send a response to the client, this is done here
		if len(minitelResponse) > 0 {
			if _, err = conn.Write([]byte(minitelResponse)); err != nil {
				logger.LogError.Print(err)
			}
		}

		if minitelParser.MinitelState == server.MINITEL_connected {
			logger.LogInfo.Print("Minitel terminal, configuring Antiope support.")
			settings.Server.Antiope = true

			// If we subsequently make a connection to another service etc. we need to relay the contents of
			//  minitelParser.Buffer to the service.
			initBytes = minitelParser.Buffer
		}

		if inputByte == 0x0d || inputByte == 0x0a || inputByte == 0x5f || inputByte == 0x23 { //CR or #

			//parse and execute the command
			cmd = parseCommand(commandBuffer)

			// ctx is only needed so that executeCommand can pass it to render()
			ctx, cancel = context.WithCancel(context.Background())
			executeCommand(ctx, conn, cmd, settings, hosts, baudRate, initBytes)

		} else {

			//validate char
			if isAlphaNumeric(inputByte) { // note that this is 7 bit only at the moment
				//add to buffer

				if inputByte == 0x8 || inputByte == 0x7f {

					// backspace
					if len(commandBuffer) > 0 {
						commandBuffer = commandBuffer[:len(commandBuffer)-1]
						writeString(conn, "\x08\x20\x08", baudRate)
					}

				} else {
					// echo back to user
					commandBuffer += string(inputByte)
					writeString(conn, string([]byte{inputByte}), baudRate)
				}

			} else {
				logger.LogInfo.Printf("Invalid character received [%xh], ignoring.", inputByte)
			}
		}

	}
}

func closeConn(conn net.Conn) {
	defer conn.Close()
	logger.LogInfo.Print("Closing connection!")
}

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

	// get a byte
	c := bufio.NewReader(conn)
	inputByte, err := c.ReadByte()

	if err != nil {
		if err != io.EOF {
			logger.LogWarn.Println("read error:", err)
		} else {
			logger.LogInfo.Println("EOF")
		}
		result = false
	}
	return result, inputByte
}
*/

func isAlphaNumeric(char byte) bool {

	// returns the ordinal value of the first character of a string
	// returns true if an ascii value for 32 - 127 and BS
	//var ord = int([]rune(char)[0])
	return (char >= 0x20 && char <= 0x7F) || char == 0x08

}

func executeCommand(ctx context.Context, conn net.Conn, cmd Command, settings config.Config, hosts map[string]string, baudRate int, initBytes []byte) {

	if cmd.isValid {

		switch cmd.name {

		case "CALL":
			url, ok := hosts[cmd.arg1]
			if !ok {
				url = cmd.arg1
			}

			// returns a bool
			_ = netClient.Connect(conn, url, settings.Pad.DLE, baudRate, initBytes)

		case "HELP":
			go render(ctx, conn, HELP, baudRate)
		case "HELPPARS":
			go render(ctx, conn, HELP_PARS, baudRate)
		case "HELPPROFILE":
			go render(ctx, conn, HELPPROFILE, baudRate)
		case "HELPCALL":
			go render(ctx, conn, HELPCALL, baudRate)
		case "HOSTS":
			go render(ctx, conn, formatHosts(hosts), baudRate)
		case "PROFILES":
			go render(ctx, conn, formatProfiles(profiles), baudRate)
		case "PARS":
			// shows current profile parameters
			go render(ctx, conn, formatPars(currentProfile), baudRate)

			// these are not go routines so we can write a prompt
		case "PROFILE":
			if len(cmd.arg1) > 0 {
				tmpProfile, ok := getProfile(cmd.arg1)
				if ok {
					go render(ctx, conn, formatProfile(tmpProfile), baudRate)
				}

			} else {
				//writeString(conn, formatProfile(currentProfile.name))
				go render(ctx, conn, fmt.Sprintf("\fCURRENT PROFILE: %s\r\n\r\n", currentProfile.name), baudRate)
			}

		case "PROFILESET":

			if setProfile(cmd.arg1) {
				logger.LogInfo.Printf("Profile set to '%s'.", cmd.arg1)
				go render(ctx, conn, OK, baudRate)
			} else {
				logger.LogError.Printf("Unable to set profile to '%s'.", cmd.arg1)
				go render(ctx, conn, ERROR, baudRate)
			}

		case "P":

			// updates the current profile
			index, err1 := strconv.Atoi(cmd.arg1)
			value, err2 := strconv.Atoi(cmd.arg2)

			if err1 != nil || err2 != nil {
				logger.LogError.Printf("Unable to set parameter '%s' to '%s'.", cmd.arg1, cmd.arg2)
				go render(ctx, conn, ERROR, baudRate)

			} else {
				if setProfileValue(&currentProfile, index, value) {
					logger.LogInfo.Printf("Parameter '%s' set to '%s'.", cmd.arg1, cmd.arg2)
					go render(ctx, conn, OK, baudRate)
					currentProfile.name = CUSTOM
				} else {
					logger.LogError.Printf("Unable to set parameter '%s' to '%s'.", cmd.arg1, cmd.arg2)
					go render(ctx, conn, ERROR, baudRate)
				}
			}
		}
	} else {
		go render(ctx, conn, HUH, baudRate)
	}
	// this is done immediately after the go routines above have STARTED not after they have finished
	//clear the buffer
	commandBuffer = ""

}
