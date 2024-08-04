package server

import (
	"github.com/johnnewcombe/telstar-library/logger"
)

const (
	MINITEL_undefined = iota
	DC_found          // 1
)

type MinitelParser struct {
	minitelState      int
	lastByteReceived  byte
	Buffer            []byte
	MinitelConnection bool
}

func (parser *MinitelParser) ParseMinitel(char byte) (byte, string) {

	// The minitel negotiation starts with DC (0x13) followed by a char in the range 40h and 5Fh
	// Mute echo until all of the negotiation is done
	var response string

	// if we have already determined a telnet connection, we dont need further parsing
	if parser.MinitelConnection {
		return char, response
	}

	if parser.minitelState == MINITEL_undefined && char == 0x13 { // DC

		logger.LogInfo.Print("DC Found")
		parser.minitelState = DC_found
		parser.Buffer = append(parser.Buffer, char)

	} else if parser.minitelState == DC_found && parser.lastByteReceived == 0x13 {

		if char >= 0x40 && char <= 0x5f {
			logger.LogInfo.Print("Minitel terminal detected.")
			parser.Buffer = append(parser.Buffer, char)
			parser.MinitelConnection = true
		}
	}

	parser.lastByteReceived = char

	if parser.minitelState != MINITEL_undefined {
		char = 0
	}

	return char, response
}
