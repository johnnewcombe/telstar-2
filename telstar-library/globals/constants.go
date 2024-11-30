package globals

const (
	ROWS = 24
	COLS = 40
)

const (
	VERSION_FILE = "version.txt"
)

const (
	GUEST_USER     = "7777777777"
	GUEST_PASSWORD = "7777"
	GUEST_NAME     = "HONOURED GUEST"
	STD_USER       = "1010101010"
	STD_PASSWORD   = "1010"
)

const (
	// delay factor
	BAUD_MAX = 0.0
	//BAUD_4800 = 0.002083 // seconds
	//BAUD_2400 = 0.004167 // seconds
	//BAUD_1200 = 0.008333 // seconds
	BAUD_4800 = 2083000 // nanoseconds
	BAUD_2400 = 4167000 // nanoseconds
	BAUD_1200 = 8333000 // nanoseconds
	BAUD_RATE = BAUD_1200

	// HASH_GUARD_TIME if a hash is received within n ms of the previous character and we are in immediate mode
	// ignore it. This should help packet radio operation where someone presses a menu choice e.g. 2
	// and hits Return (hash) as the TNC will send 2# in this scenario the hash needs to be ignored.
	HASH_GUARD_TIME    = 200
	CONNECT_DELAY_SECS = 2
)

const (
	DBPRIMARY   = "primary"
	DBSECONDARY = "secondary"
)
const (
	NULL     byte = 0x00
	BS       byte = 0x08
	HT       byte = 0x09
	LF       byte = 0x0A
	VT       byte = 0x0B
	CLS      byte = 0x0C
	CR       byte = 0x0D
	DC       byte = 0x13
	SUB      byte = 0x1A
	ESC      byte = 0x1B
	HOME     byte = 0x1E
	SPC      byte = 0x20
	CURON    byte = 0x11
	CUROFF   byte = 0x14
	SO       byte = 0x0e
	SI       byte = 0x0f
	US       byte = 0x1f
	HASH     byte = 0x5F
	POUND    byte = 0x23
	ASTERISK byte = 0x2A
	DEL      byte = 0x7F
)

const (
	SEPARATOR_GRAPHIC_DOTS_LOW     = "000000000000000000000000000000000000000"
	SEPARATOR_GRAPHIC_DOTS_MID     = "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$"
	SEPARATOR_GRAPHIC_DOTS_HIGH    = "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
	SEPARATOR_GRAPHIC_SOLID_LOW    = "ppppppppppppppppppppppppppppppppppppppp"
	SEPARATOR_GRAPHIC_SOLID_MID    = ",,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,"
	SEPARATOR_GRAPHIC_SOLID_HIGH   = "#######################################"
	SEPARATOR_GRAPHIC_SOLID_DOUBLE = "sssssssssssssssssssssssssssssssssssssss"
)

const (
	ALPHA_BLACK    = "\x1b\x40"
	ALPHA_RED      = "\x1b\x41"
	ALPHA_GREEN    = "\x1b\x42"
	ALPHA_YELLOW   = "\x1b\x43"
	ALPHA_BLUE     = "\x1b\x44"
	ALPHA_MAGENTA  = "\x1b\x45"
	ALPHA_CYAN     = "\x1b\x46"
	ALPHA_WHITE    = "\x1b\x47"
	FLASH          = "\x1b\x48"
	STEADY         = "\x1b\x49"
	ENDBOX         = "\x1b\x4a" // Antiope only
	STARTBOX       = "\x1b\x4b" // Antiope only
	NORMAL_HEIGHT  = "\x1b\x4c"
	DOUBLE_HEIGHT  = "\x1b\x4d"
	DOUBLE_WIDTH   = "\x1b\x4e" // Antiope only
	DOUBLE_SIZE    = "\x1b\x4f" // Antiope only
	MOSAIC_BLACK   = "\x1b\x50"
	MOSAIC_RED     = "\x1b\x51"
	MOSAIC_GREEN   = "\x1b\x52"
	MOSAIC_YELLOW  = "\x1b\x53"
	MOSAIC_BLUE    = "\x1b\x54"
	MOSAIC_MAGENTA = "\x1b\x55"
	MOSAIC_CYAN    = "\x1b\x56"
	MOSAIC_WHITE   = "\x1b\x57"
	CONCEAL        = "\x1b\x58"
	STOP_LINING    = "\x1b\x59" // Antiope only
	START_LINING   = "\x1b\x5a" // Antiope only
	CSI            = "\x1b\x5b" // Antiope only
	END_BACKGROUND = "\x1b\x5c"
	NEW_BACKGROUND = "\x1b\x5d"
	HOLD_MOSAIC    = "\x1b\x5e"

	RELEASE_MOSAIC = "\x1b\x5f"
)

const MINITEL_ENQ_ROM = "\x0c\x1e\x1e\x1e\x1b\x39\x7b"

const (
	FRAME_TYPE_INITIAL      = "initial"
	FRAME_TYPE_MAININDEX    = "mainindex"
	FRAME_TYPE_INFORMATION  = "information"
	FRAME_TYPE_EXIT         = "exit"
	FRAME_TYPE_GATEWAY      = "gateway"
	FRAME_TYPE_TEST         = "test"
	FRAME_TYPE_TELESOFTWARE = "telesoftware"
	FRAME_TYPE_RESPONSE     = "response"
	FRAME_TYPE_SYSTEM       = "system"
	FRAME_TYPE_EXCEPTION    = "exception"
)

const (
	CONTENT_TYPE_EDITTF = "edit.tf"
	CONTENT_TYPE_MARKUP = "markup"
	CONTENT_TYPE_RAW    = "raw"
	CONTENT_TYPE_RAWV   = "rawV"
	CONTENT_TYPE_RAWT   = "rawT"
	CONTENT_TYPE_ZXNET  = "zxnet"
)

const (
	ConnectionModeViewdata       = "viewdata"         // Full duplex Viewdata service, this will not need buffering.
	ConnectionModeFullDuplex     = "full_duplex"      // Full duplex where gateway service does not echo the data.
	ConnectionModeFullDuplexEcho = "full_duplex_echo" // Full duplex where gateway service echos the data.
	ConnectionModeHalfDuplex     = "half_duplex"      // Line mode where data is terminated by CR and responds in a similar manner. no echo from server

)

// This is the telstar logo, defined here to remove dependencies on external files
const TELSTAR_LOGO = "                  [G]T\r\n                [G]T[R]E[C]L\r\n              [G]T[R]E[C]L[B]S[W]T\r\n            [G]T[R]E[C]L[B]S[W]T[M]A[Y]R\r\n              [G]L[R]S[C]T[B]A[W]R\r\n                [G]T[R]A[C]R\r\n                  [G]R\r\n"
const EDITTF_LOGO = "http://edit.tf/#0:QIECBAgQIECBAgQIECBAgQAqiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIEAKoBihpiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBACqAYoaYEph6iBAgQIECBAgQIECBAgQIECBAgQIECBAgQAqgGKGmBKYeoFgg6SBAgQIECBAgQIECBAgQIECBAgQIECBAgQBpgSmHqBYIOkgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQB6gWCDpIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQA6SBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"

/*
   "logo":
       [
           "http://edit.tf/#0:QIECBAgQIECBAgQIECBAgQAqiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIEAKoBihpiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBACqAYoaYEph6iBAgQIECBAgQIECBAgQIECBAgQIECBAgQAqgGKGmBKYeoFgg6SBAgQIECBAgQIECBAgQIECBAgQIECBAgQBpgSmHqBYIOkgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQB6gWCDpIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQA6SBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA",
           0, 7, 0, 39
       ],
*/
