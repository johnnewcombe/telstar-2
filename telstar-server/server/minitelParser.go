package server

import (
	"github.com/johnnewcombe/telstar-library/logger"
)

const (
	MINITEL_undefined = iota
	MINITEL_DC_found  // 1
	MINITEL_ENQ_ROM_start_found
	MINITEL_vendor_found
	MINITEL_model_found
	MINITEL_revision_found
	MINITEL_not_connected
	MINITEL_connected
)

type MinitelParser struct {
	MinitelState     int
	lastByteReceived byte
	Buffer           []byte
	Vendor           byte
	Model            byte
	Revision         byte
}

func (parser *MinitelParser) ParseMinitelEnqRom(char byte) (byte, string) {

	// use the minitel parser later or wait here or move it to here
	// result should be 01/XX/YY/ZZ/04 where :

	// XX is the vendor code (41 = Matra, 42 = Philips, 43 = Alcatel, and more exotic ones after)
	// YY is the model code (many exist – always printable characters)
	// ZZ is the revision code (many exist – always printable characters)

	// TODO: Another method needed, perhaps *MINITEL
	// Note that very early Minitels (ABCDEF keyboards) did respond to this
	// ENQ_ROM sequence only if it came from the modem.

	var response string

	// FIXME: Complete this parser

	// if we have already determined a minitel connection, we don't need further parsing
	if parser.MinitelState == MINITEL_connected ||
		parser.MinitelState == MINITEL_not_connected {
		return char, response
	}

	if parser.MinitelState == MINITEL_undefined && char&0x7f == 0x01 { // SOH
		logger.LogInfo.Print("Minitel ENQ_ROM start received")
		parser.MinitelState = MINITEL_ENQ_ROM_start_found
		parser.Buffer = append(parser.Buffer, char)

	} else if parser.MinitelState == MINITEL_ENQ_ROM_start_found &&
		char&0x7f >= 0x0 && char <= 0x7f {
		logger.LogInfo.Print("Minitel ENQ_ROM vendor received")
		parser.MinitelState = MINITEL_vendor_found
		parser.Vendor = char & 0x7f

	} else if parser.MinitelState == MINITEL_vendor_found &&
		char&0x7f >= 0x0 && char <= 0x7f {
		logger.LogInfo.Print("Minitel ENQ_ROM model received")
		parser.MinitelState = MINITEL_model_found
		parser.Model = char & 0x7f

	} else if parser.MinitelState == MINITEL_model_found &&
		char&0x7f >= 0x20 && char <= 0x7f {
		logger.LogInfo.Print("Minitel ENQ_ROM revision received")
		parser.MinitelState = MINITEL_revision_found
		parser.Revision = char & 0x7f

	} else if parser.MinitelState == MINITEL_revision_found && char&0x7f == 0x04 { //EOT
		logger.LogInfo.Print("Minitel ENQ_ROM end received")
		parser.MinitelState = MINITEL_connected

	} else {
		parser.MinitelState = MINITEL_not_connected
	}

	parser.Buffer = append(parser.Buffer, char)
	parser.lastByteReceived = char

	if parser.MinitelState != MINITEL_undefined {
		char = 0
	}

	return char, response
}

func (parser *MinitelParser) ParseMinitelDc(char byte) (byte, string) {

	// The minitel negotiation starts with DC (0x13) followed by a char in the range 40h and 5Fh
	// Mute echo until all of the negotiation is done
	var response string

	// if we have already determined a minitel connection, we dont need further parsing
	if parser.MinitelState == MINITEL_connected {
		return char, response
	}

	if parser.MinitelState == MINITEL_undefined && char == 0x13 { // DC

		logger.LogInfo.Print("DC Found")
		parser.MinitelState = MINITEL_DC_found
		parser.Buffer = append(parser.Buffer, char)

	} else if parser.MinitelState == MINITEL_DC_found && parser.lastByteReceived == 0x13 {

		if char >= 0x40 && char <= 0x5f {
			logger.LogInfo.Print("Minitel terminal detected.")
			parser.Buffer = append(parser.Buffer, char)
			parser.MinitelState = MINITEL_connected
		}
	}

	parser.lastByteReceived = char

	if parser.MinitelState != MINITEL_undefined {
		char = 0
	}

	return char, response
}
