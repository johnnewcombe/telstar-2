package constants

const (
	ROWS             = 24
	COLS             = 40
	//SIZE             = 1.0 //works to about 1.5 after that we see gaps in the graphics
	//TEXTSIZE         = 23
	//GRAPHICSSIZE     = 20.295 * SIZE
	//TEXTWIDTH        = 16 * SIZE
	//TEXTHEIGHT       = 20 * SIZE
	FLASHPERIOD_MS   = 250
	CURSOR_THICKNESS = 3.0
)

const (
	Padding = 5

	// Width and Height
	//ScreenWidth   = COLS * TEXTWIDTH
	//ScreenHeight  = ROWS * TEXTHEIGHT
	//ToolBarWidth  = ScreenWidth
	//ToolBarHeight = 22
	//StatusWidth   = ScreenWidth
	//StatusHeight  = 22
	//WindowWidth   = ScreenWidth + Padding*2
	//WindowHeight  = ScreenHeight + ToolBarHeight + StatusHeight + Padding*4

	// X and Y
	//LeftMargin = 0//250
	//ToolBarX = LeftMargin
	//ToolBarY = Padding // small amount of padding
	//ScreenX  = ToolBarX
	//ScreenY  = ToolBarY + ToolBarHeight + Padding
	//StatusX  = 10
	//StatusY  = ScreenY + ScreenHeight + Padding
)

const (
	ALPHA_RED           byte = 0x41
	ALPHA_GREEN         byte = 0x42
	ALPHA_YELLOW        byte = 0x43
	ALPHA_BLUE          byte = 0x44
	ALPHA_MAGENTA       byte = 0x45
	ALPHA_CYAN          byte = 0x46
	ALPHA_WHITE         byte = 0x47
	FLASH               byte = 0x48
	STEADY              byte = 0x49
	NORMAL_HEIGHT       byte = 0x4c
	DOUBLE_HEIGHT       byte = 0x4d
	MOSAIC_RED          byte = 0x51
	MOSAIC_GREEN        byte = 0x52
	MOSAIC_YELLOW       byte = 0x53
	MOSAIC_BLUE         byte = 0x54
	MOSAIC_MAGENTA      byte = 0x55
	MOSAIC_CYAN         byte = 0x56
	MOSAIC_WHITE        byte = 0x57
	CONCEAL             byte = 0x58
	CONTIGUOUS_GRAPHICS byte = 0x59
	SEPARATED_GRAPHICS  byte = 0x5a
	BLACK_BACKGROUND    byte = 0x5c
	NEW_BACKGROUND      byte = 0x5d
	HOLD_GRAPHICS       byte = 0x5e
	RELEASE_GRAPHICS    byte = 0x5f
)

const (
	NULL     byte = 0x00
	BS       byte = 0x08
	HT       byte = 0x09
	LF       byte = 0x0A
	VT       byte = 0x0B
	CLEAR    byte = 0x0C
	CR       byte = 0x0D
	ESC      byte = 0x1B
	HOME     byte = 0x1E
	CURON    byte = 0x11
	CUROFF   byte = 0x14
	SPACE    byte = 0x20
	POUND    byte = 0x23
	ASTERISK byte = 0x2A
	HASH     byte = 0x5F
)
