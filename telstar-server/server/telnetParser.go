package server

import (
	"github.com/johnnewcombe/telstar-library/logger"
)

const (
	TELNET_undefined = iota
	IAC_found        // 1
	SB_found         // 2
	DO_found         // 3
	WILL_found       // 4
	DONT_found       // 5
	WONT_found       // 6
)

type TelnetParser struct {
	telnetState      int
	lastByteReceived byte
	suppressGoAhead  bool
	TelnetConnection bool
}

/*
   If you send an IAC WILL TRANSMIT-BINARY, you are talking about the direction from you to the peer and you are
   effectively asking the peer to allow you send binary. The peer can answer either  DO or DONT.

   If you send an IAC DO TRANSMIT-BINARY, you are instead talking about the direction from the peer to you and you
   are effectively asking the peer to start transmitting binary if it is able to and willing. The peer can answer
   either WILL or WONT.

   However, if the peer sends an IAC WILL TRANSMIT-BINARY (as opposed to you sending it), it is talking about the
   direction from the peer to you and it is effectively asking you to allow it send binary. You can answer
   either DO or DONT.

   Analogously, if the peer sends an IAC DO TRANSMIT-BINARY, it is talking about the direction from you to the
   peer and it is effectively asking you to start transmitting binary if you are able to and willing. You can answer
   either WILL or WONT.
*/

func (parser *TelnetParser) ParseTelnet(char byte) (byte, string) {

	// The telnet negotiation starts with IAC (FF) followed by 2 bytes unless the 2nd byte is SB (FA) Sub-Negotiation
	// in which case all subsequent chars can be ignored until the sequence SE (FF FO) Sub-Negotiation End
	// Mute echo until all of the negotiation is done
	var response string

	// if we have already determined a telnet connection, we dont need further parsing
	if parser.TelnetConnection {
		return char, response
	}

	// State machine to sort out telnet parsing
	if parser.telnetState == TELNET_undefined && char == 0xff { // IAC

		logger.LogInfo.Print("IAC Start")
		parser.telnetState = IAC_found

	} else if parser.telnetState == IAC_found && char == 0xff && parser.lastByteReceived == 0xff {

		logger.LogInfo.Print("IAC Cancelled for FF FF")
		parser.telnetState = TELNET_undefined
		char = 0 // because we are setting to undefined, we need to explicitly set char to 0

	} else if parser.telnetState == IAC_found && char == 0xfa { // SB. sub negotiation start

		logger.LogInfo.Print("SB Started")
		parser.telnetState = SB_found

	} else if parser.telnetState == SB_found && char == 0xf0 { // SE, end of Sub-Negotiation

		// end of the command (sub negotiation)
		logger.LogInfo.Print("SB Ended")
		parser.telnetState = TELNET_undefined
		char = 0 // because we are setting to undefined, we need to explicitly set char to 0

	} else if parser.telnetState == IAC_found {

		if char >= 240 && char <= 255 {
			logger.LogInfo.Print("Telnet client detected.")
			parser.TelnetConnection = true
		}

		// These will all be followed by options
		if char == 251 { // WILL
			parser.telnetState = WILL_found
		} else if char == 252 { // WONT
			parser.telnetState = WONT_found
		} else if char == 253 { // DO
			parser.telnetState = DO_found
		} else if char == 254 { // DONT
			parser.telnetState = DONT_found
		} else if char >= 240 && char <= 250 { // ignore these
			// no options for these so telnet complete
			logger.LogInfo.Printf("% x command ignored IAC End", char)
			parser.telnetState = TELNET_undefined
			char = 0 // because we are setting to undefined, we need to explicitly set char to 0
		}

	} else if parser.telnetState == WILL_found {

		// anything arriving here is an option
		logger.LogInfo.Print("WILL % x", char)
		parser.telnetState = TELNET_undefined
		char = 0 // because we are setting to undefined, we need to explicitly set char to 0

	} else if parser.telnetState == WONT_found {

		// anything arriving here is an option
		logger.LogInfo.Print("WONT % x", char)
		parser.telnetState = TELNET_undefined
		char = 0 // because we are setting to undefined, we need to explicitly set char to 0

	} else if parser.telnetState == DO_found {

		// anything arriving here is an option
		logger.LogInfo.Print("DO % x", char)

		if char == 0x01 { // ECHO
			// return WILL echo in response to receiving the DO
			response = "\xff\xfb\x01"
		} else if char == 0x03 { // GOAHEAD
			// return WILL suppress go ahead in response to receiving the DO
			// also indicate that I am will be echoing by sending WILL echo
			response = "\xff\xfb\x03"
		} else {
			// send WONT to all other do requests
			response = "\xff\xfc"
		}

		logger.LogInfo.Print("IAC End")
		parser.telnetState = TELNET_undefined
		char = 0 // because we are setting to undefined, we need to explicitly set char to 0

	} else if parser.telnetState == DONT_found {

		// anything arriving here is an option
		logger.LogInfo.Print("DONT  % x", char)
		parser.telnetState = TELNET_undefined
		char = 0 // because we are setting to undefined, we need to explicitly set char to 0

	}

	// Telnet clients will either send CR/LF or CR/NULL the NULL can cause issues when detecting EOL
	// as it tends to around at the start of subsequent input strings
	if char == 0x0a || char == 0x0d {
		char = 0
	}

	parser.lastByteReceived = char

	// any state other than undefined will result in the char being ignored
	if parser.telnetState != TELNET_undefined {
		char = 0
	}

	return char, response
}
