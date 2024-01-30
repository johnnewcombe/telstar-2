package display

import (
	"bitbucket.org/johnnewcombe/telstar-client/constants"
	"bitbucket.org/johnnewcombe/telstar-client/customFyne"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"sync"
	"time"
)

var colourMap = make(map[int]color.Color)

// Write should be concurrent safe. Writing updates a cache and the refresh then displays that cache.
// Documentation for fyne 'suggests' that this is concurrent safe. Writing 0 will cause the attributes
// to be written but not the text. This allows background or flash etc to be set without changing the
// text.
func (s *Screen) Write(char byte) error {

	// At this point the cursor could have been positioned to an invalid column or row
	// e.g. row 24 or col 40. This can only be fixed here now that we know what the next char is.
	done := s.processControls(char)

	if s.cursorCol > constants.COLS-1 || s.cursorRow > constants.ROWS-1 || s.cursorCol < 0 || s.cursorRow < 0 {
		return fmt.Errorf("bad index col: %d, row: %d", s.cursorCol, s.cursorRow)
	}

	// cursor could have changed
	s.setCursor(s.getIndex(s.cursorCol, s.cursorRow))

	// if handled above then we are done
	if !done {
		// if not ctrl then muct be a char to display
		s.write(char)
		s.cursorCol++
		s.wrap()
		s.roll()
		s.setCursor(s.getIndex(s.cursorCol, s.cursorRow))
	}

	return nil
}

func (s *Screen) write(char byte) {

	// do not process a lower row of a double height row pair
	if s.isRowDisabled(s.cursorRow) {
		return
	}

	var ichar int

	// if escape has already been set then we are dealing with a control character
	if s.escape {

		// add the control controlValue to the attributes
		s.attributes[s.getIndex(s.cursorCol, s.cursorRow)].controlValue = char
		// recalculate the attributes from the previous (or default) to this cell and to the end of the row
		s.setAttributes(s.cursorCol, s.cursorRow)

		// get the current attribute after all of the changes
		attrib := s.attributes[s.getIndex(s.cursorCol, s.cursorRow)]

		if char == constants.RELEASE_GRAPHICS {
			print("stop")
		}

		if s.holdGraphicsMode && attrib.graphics {

			// Hold graphics will show the previously stored graphic within that frame irrespective of the
			// line it is on instead of a blank for where the control code would have been.
			// There are a few things to note
			// 1 if the control code is a graphic colour change, the hold graphic retains the graphic
			//   colour that was in use before the graphic colour change code. Everything after the graphic
			//   colour change code will take on the new colour.
			// 2 the initial hold graphics code is displayed as a blank, subsequent hold graphics codes are
			//   displayed using the stored (held) graphic.
			// 3 The Release Graphic code is replaced by the hold graphic.
			// 4 If there is no graphic has been displayed when the Hold Graphic occurs, a blank is displayed for
			//   the control code as usual until a graphic character is displayed.
			ichar = s.heldGraphic.character

		} else {
			// display a space in place of the control character
			ichar = 0x20
		}

		// setDoubleHeightText is the same as setText but handles the lower row also
		// and will get called for either a hold graphics or if it is the DH control
		// code, a space will be output
		if s.isDoubleHeightRow(s.cursorRow) {
			s.setDoubleHeightText(ichar, s.cursorCol, s.cursorRow, attrib.conceal && !s.revealMode)
		} else {
			if ichar != 0 {
				s.setText(ichar, s.cursorCol, s.cursorRow, attrib.conceal && !s.revealMode, customFyne.CharTypeNone)
			}
		}

		// set any global display settings
		switch char {
		case constants.HOLD_GRAPHICS:
			s.holdGraphicsMode = true
		case constants.RELEASE_GRAPHICS:
			// we need to delay setting the s.holdGraphicsMode=false as we may
			// still have to process the releaseGraphic character itself(see below)
			s.holdGraphicsMode = false
		}
		s.escape = false

	} else {

		// not an escaped char so update the attributes and add directly to the grid
		s.attributes[s.getIndex(s.cursorCol, s.cursorRow)].controlValue = 0

		// recalculate the attributes from the previous (or default) to this cell and to the end of the row
		s.setAttributes(s.cursorCol, s.cursorRow)

		attrib := s.attributes[s.getIndex(s.cursorCol, s.cursorRow)]
		ichar = int(char)

		// check for graphics
		if !attrib.graphics {
			// fix the hash and other odd symbols with a map
			mchar := s.alphaCharMap[ichar]
			if mchar != 0 {
				ichar = mchar
			}

		} else {

			// save to use as held graphic
			if s.holdGraphicsMode {
				s.heldGraphic.character = ichar
				s.heldGraphic.attributes = attrib
			}
		}

		// double height
		if s.isDoubleHeightRow(s.cursorRow) {
			s.setDoubleHeightText(ichar, s.cursorCol, s.cursorRow, attrib.conceal && !s.revealMode)
		} else {
			s.setText(ichar, s.cursorCol, s.cursorRow, attrib.conceal && !s.revealMode, customFyne.CharTypeNone)
		}
	}
}

// processControls returns true if a control character was processed, false otherwise.
func (s *Screen) processControls(char byte) bool {

	// check for normal CTRL Codes, these represent things that affect the whole screen not just a
	// single text cell
	processed := true

	// fix things up before processing
	s.wrap()
	s.roll()

	// we dont need to worry about wrap and roll for these
	switch char {
	case constants.CUROFF:
		s.cursorOn = false
	case constants.CURON:
		s.cursorOn = true
	case constants.ESC:
		s.escape = true

	default:

		// these are affected by wrap and roll so
		switch char {
		case constants.CLEAR:
			s.Clear()
		case constants.HOME:
			s.home()
		case constants.CR:
			s.cursorCol = 0
		case constants.VT:
			s.cursorRow--
		case constants.BS:
			s.cursorCol--
		case constants.LF:
			s.cursorRow++
		case constants.HT:
			s.cursorCol++
		default:
			processed = false
		}

		// fix things up after processing controls
		if processed {
			s.wrap()
			s.roll()
		}
	}
	return processed
}

func (s *Screen) setText(char int, cursorCol int, cursorRow int, hidden bool, charType customFyne.CharType) {

	index := s.getIndex(cursorCol, cursorRow)

	// ignore if a Double Height lower row
	if s.isRowDisabled(cursorRow) && !s.attributes[index].doubleHeight && char != 0 {
		return
	}

	// reset the modes to defaults
	if cursorCol == 0 {
		s.holdGraphicsMode = false
		s.doubleHeightMode = false
	}

	if char > 0 {

		if s.attributes[index].graphics {

			if !s.isBlastThrough(char) {

				if s.Debug {
					s.debugWrite(byte(char), "Graphics")
				}

				// Hold Graphics is a little convoluted. See note within the write() function.
				mosaic := s.getMosaicContainer(index)

				// FIXME set block type i.e. DH Upper/Lower or Normal.
				if s.holdGraphicsMode {
					mosaic.SetMosaic(byte(char), s.attributes[index].foreColour, s.heldGraphic.attributes.nonContiguous, charType)
				} else {
					mosaic.SetMosaic(byte(char), s.attributes[index].foreColour, s.attributes[index].nonContiguous, charType)
				}
				mosaic.Hidden = hidden
				mosaic.Refresh()

			} else {

				if s.Debug {
					s.debugWrite(byte(char), "Blast Through")
				}

				// two chars in the set 40-5f need substitutions when in graphic mode
				// so that blast through graphics show the hash and half symbol correctly.
				// Note that isBlastThrough() will return false after any substitution as
				// 0x23 and  0xbd are not valid blast through characters.
				if s.isBlastThrough(char) {
					if char == 0x5f { // hash
						char = 0x23
					}
					if char == 0x5c { // fraction 1/2
						char = 0xbd
					}
				}

				txt := s.getTextCanvas(index)
				txt.Text = string(rune(char))
				txt.Hidden = hidden
				txt.Color = s.attributes[index].foreColour
				txt.Refresh()
			}

		} else {

			if s.Debug {
				s.debugWrite(byte(char), "Alphanumeric")
			}

			txt := s.getTextCanvas(index)
			txt.Text = string(rune(char))
			txt.Hidden = hidden
			txt.Color = s.attributes[index].foreColour
			txt.Refresh()

		}
	}

	bg := s.getBackgroundCanvas(index)
	bg.FillColor = s.attributes[index].backColour
	bg.Hidden = hidden
	bg.Refresh()

}

func (s *Screen) setDoubleHeightText(char int, cursorCol int, cursorRow int, hidden bool) {

	var upper, lower byte

	attr1 := s.attributes[s.getIndex(cursorCol, cursorRow)]
	blastThough := s.isBlastThrough(char)

	// this routine is looking at the upper row and both double height chars
	// and none double height chars (i.e. after a Normal Height Code) could
	// appear here
	if attr1.doubleHeight && attr1.controlValue == 0 {
		// set the upper part this will apply the colours etc
		if attr1.graphics && !blastThough {
			upper, lower = s.getDoubleHeightMosaic(byte(char))
			char = int(upper)
		} else {
			// char is alphanumeric, e000 is start of double height (top row) characters
			char += 0xe000
		}
		s.setText(char, cursorCol, cursorRow, hidden, customFyne.CharTypeDoubleHeightUpper)
	} else {
		// on the top row of a double height pair but not double height
		// e.g. appears after a Normal Text control so
		// just display normal text
		s.setText(char, cursorCol, cursorRow, hidden, customFyne.CharTypeNone)
	}

	// if we are on the last row then nothing more to do
	if s.isLastRow(cursorRow) {
		return
	}

	// both double height chars and none double height chars could cause us to
	// get here, if they are non-double height then only the attributes need
	// to be set as this routine is looking at the lower row.
	if attr1.doubleHeight && attr1.controlValue == 0 {
		// display the lower part of the char
		if attr1.graphics && !blastThough {
			//_, lower := s.getDoubleHeightMosaic(byte(char))
			char = int(lower)
		} else {
			char += 0x100
		}
		s.setText(char, cursorCol, cursorRow+1, hidden, customFyne.CharTypeDoubleHeightLower)
	}

}

func (s *Screen) removeMosaicOffset(b byte) byte {

	if b >= 0x60 && b <= 0x7f {
		return b & 0b10111111
	} else if b >= 0x20 && b <= 0x3f {
		return b & 0b11011111
	} else {
		return b
	}
}

func (s *Screen) addMosaicOffset(b byte) byte {

	if b >= 0x00 && b <= 0x1f {
		return b | 0b00100000
	} else if b >= 0x20 && b <= 0x3f {
		return b | 0b01000000
	} else {
		return b
	}

}

//getDoubleHeightMosaicUpper Accepts mosaic char and returns upper and lower portions of double height chars
func (s *Screen) getDoubleHeightMosaic(c byte) (upper byte, lower byte) {

	c = s.removeMosaicOffset(c)

	dhCharU := func(char byte) byte {

		var result, b, mask byte

		mask = 0b00000001

		for i := 0; i < 4; i++ {
			b = char & mask // get bit
			if i < 2 {
				result = result | b // set the bit
			}
			b = b << 2
			result = result | b // set the bit + 2

			mask = mask << 1

		}

		return result
	}
	dhCharL := func(char byte) byte {

		var result, b, mask byte

		mask = 0b00100000

		for i := 0; i < 4; i++ {
			b = char & mask // get bit
			if i < 2 {
				result = result | b // set the bit
			}
			b = b >> 2
			result = result | b // set the bit + 2

			mask = mask >> 1
		}

		return result
	}
	upper = dhCharU(c)
	//c = c >> 3
	lower = dhCharL(c)
	//c = c << 3

	return s.addMosaicOffset(upper), s.addMosaicOffset(lower)
}

// setAttributes runs throught all attributes from the specified position to the end
// of the row. The function also syncronises the lower row of a double height row pair
// with the upper row.
func (s *Screen) setAttributes(startCol int, row int) {

	var (
		byt           byte
		currentAttrib Attributes
	)
	if s.isDoubleHeightLowerRow(row) {
		return
	}

	// loop though each cell in the row from the next cell to the end of the row
	for i := s.getIndex(startCol, row); i < s.getRowEndIndex(row); i++ {

		// attributes will never be applied to the first column but any control value
		// that is present will be used to affect following cols
		// this is the case all along in that the control code of the current column
		// is used to change the following cols

		// get the current attrib for the position
		currentAttrib = s.attributes[i]
		byt = currentAttrib.controlValue

		// update the next cell with the current
		currentAttrib.copyTo(&s.attributes[i+1])

		// Some docs suggest that we need to detect a change in colour and/or Double height
		// in order to determine if the held graphic should be cleared. Rather than resetting it
		// on every colour setting. However, this may not be the case as the test card uses
		// Alpha Yellow, Conceal, <some text>, Alpha Yellow and the second alpha yellow cancels conceal

		// update the current attribute with any changes from this cell
		// mosaics
		if byt >= constants.MOSAIC_RED && byt <= constants.MOSAIC_WHITE {
			s.attributes[i+1].graphics = true
			s.attributes[i+1].conceal = false
			s.attributes[i+1].foreColour = colourMap[int(byt)-0x50]
		}
		// alpha
		if byt >= constants.ALPHA_RED && byt <= constants.ALPHA_WHITE {
			s.attributes[i+1].graphics = false
			s.attributes[i+1].conceal = false
			s.attributes[i+1].foreColour = colourMap[int(byt)-0x40]
		}

		// special cases
		switch byt {
		case constants.BLACK_BACKGROUND:
			// this is unusual in that the position where the control code is, is also affected
			s.attributes[i].backColour = colourMap[0]
			s.attributes[i+1].backColour = colourMap[0]
		case constants.DOUBLE_HEIGHT:
			s.doubleHeightMode = true
			s.attributes[i+1].doubleHeight = true
			s.attributes[i+1].conceal = false

			// This makes this row a double height row so we need to copy the attributes
			// that have been set prior to the current column again so that the row below
			// will be correct. Now that this is a double height row, future columns will
			// be handled by code below in the normal way.
			for index := s.getIndex(0, row); index < i; index++ {
				attr := s.attributes[index]
				attr.copyTo(&s.attributes[index+constants.COLS])
				c, r := s.getColRow(index + constants.COLS)
				s.setText(0, c, r, attr.conceal && !s.revealMode, customFyne.CharTypeNone)
			}
		case constants.NORMAL_HEIGHT:
			s.doubleHeightMode = false
			s.attributes[i+1].doubleHeight = false
			s.attributes[i+1].conceal = false
		case constants.NEW_BACKGROUND:
			// this is unusual in that the position where the control code is, is also affected
			s.attributes[i].backColour = s.attributes[i].foreColour
			s.attributes[i+1].backColour = s.attributes[i].foreColour
		case constants.FLASH:
			s.attributes[i+1].flash = true
		case constants.STEADY:
			s.attributes[i+1].flash = false
		case constants.CONCEAL:
			// there is no reveal code, reveal is executed by the user pressing TAB
			// normally a control code affects the following cells [i+1], however conceal
			// is unusual in that the position where the control code is, is also affected
			// this wouldn't normally matter as control codes are blank unless Hold Graphics is
			// switched on in which case the last used graphic is displayed instead.
			s.attributes[i].conceal = true   // current cell
			s.attributes[i+1].conceal = true // next cell
		case constants.CONTIGUOUS_GRAPHICS:
			s.attributes[i+1].nonContiguous = false
		case constants.SEPARATED_GRAPHICS:
			s.attributes[i+1].nonContiguous = true

		}

		// we need to write out the rest of the line in order to set the attributes for the line
		// this is important as the next char may shift focus away from this line to another.
		c, r := s.getColRow(i)
		s.setText(0, c, r, s.attributes[i].conceal && !s.revealMode, customFyne.CharTypeNone)

		// the main function is always updating the next column so index (i) never
		// gets to the very last column which means that we need to set this explicitly
		if i == s.getRowEndIndex(row)-1 {
			c, r := s.getColRow(i + 1)
			s.setText(0, c, r, s.attributes[i+1].conceal && !s.revealMode, customFyne.CharTypeNone)
		}

		// keep lower row of a double height row in sync with upper row
		if s.isDoubleHeightRow(row) && !s.isLastRow(row) {
			// copy all of the attributes from this row to the following row
			s.attributes[i].copyTo(&s.attributes[i+constants.COLS])
			c, r := s.getColRow(i + constants.COLS)
			s.setText(0, c, r, s.attributes[i].conceal && !s.revealMode, customFyne.CharTypeNone)

			// the main function is always updating the next column so index (i) never
			// gets to the very last column which means that we need to set this explicitly
			if i == s.getRowEndIndex(row)-1 {
				s.attributes[i+1].copyTo(&s.attributes[i+1+constants.COLS])
				c, r := s.getColRow(i + 1 + constants.COLS)
				s.setText(0, c, r, s.attributes[i+1].conceal && !s.revealMode, customFyne.CharTypeNone)
			}
		}

	}
}

func (s *Screen) home() {
	s.cursorRow = 0
	s.cursorCol = 0
	s.setCursor(s.getIndex(s.cursorCol, s.cursorRow))
}

func (s *Screen) Clear() {

	// reset the cursor
	s.home()
	s.revealMode = false
	s.heldGraphic.character = 0x20
	s.heldGraphic.attributes.clear()
	s.flashHidePeriod = false
	s.cursorOn = false

	for g := 0; g < constants.ROWS*constants.COLS; g++ {

		// reset all attributes to defaults
		s.attributes[g].clear()

		text := s.getTextCanvas(g)
		text.Text = " "
		text.Color = colourMap[7]
		text.Hidden = false
		text.Refresh()

		// clear the mosaic
		mosaic := s.getMosaicContainer(g)
		mosaic.Clear()

		// clear the background rectangle
		background := s.getBackgroundCanvas(g)
		background.FillColor = colourMap[0]
		background.Hidden = false
		background.Refresh()

		cursor := s.getCursorCanvas(g)
		cursor.FillColor = colourMap[7]
		cursor.Hidden = true
		cursor.StrokeWidth = 0

		cursor.Refresh()
	}
}

func (s *Screen) HoldGraphics() {
	s.holdGraphicsMode = true
}

func (s *Screen) ReleaseGraphics() {
	s.holdGraphicsMode = false
}

func (s *Screen) Conceal() {
	// conceal
	s.revealMode = false

	for i := 0; i < constants.ROWS*constants.COLS; i++ {

		text := s.getTextCanvas(i)
		text.Hidden = s.attributes[i].conceal && !s.revealMode
		text.Refresh()

		bg := s.getBackgroundCanvas(i)
		bg.Hidden = s.attributes[i].conceal && !s.revealMode
		bg.Refresh()

		mosaic := s.getMosaicContainer(i)
		mosaic.Hidden = s.attributes[i].conceal && !s.revealMode
		mosaic.Refresh()

	}
}

func (s *Screen) Reveal() {
	// reveal
	s.revealMode = true
	for i := 0; i < constants.ROWS*constants.COLS; i++ {
		text := s.getTextCanvas(i)
		text.Hidden = false
		text.Refresh()

		bg := s.getBackgroundCanvas(i)
		bg.Hidden = false
		bg.Refresh()

		mosaic := s.getMosaicContainer(i)
		mosaic.Hidden = s.attributes[i].conceal && !s.revealMode
		mosaic.Refresh()

	}
}

func (s *Screen) runFlashTimer(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		// process any requested cancellation by checking the Done channel of the context
		// the Done channel is set in the parent function in response to this function
		// indicating a DLE (Data Link Escape) on the dleSignal channel. See below
		select {
		case <-ctx.Done():
			// ctx is telling us to stop,
			//log.Println("Screen.runFlashTimer() goroutine cancelled")
			return

		default:
		}

		// toggle the flash status
		s.flashHidePeriod = !s.flashHidePeriod

		time.Sleep(time.Millisecond * constants.FLASHPERIOD_MS)
		for i := 0; i < constants.ROWS*constants.COLS; i++ {

			if s.attributes[i].cursor && s.cursorOn {
				cursor := s.getCursorCanvas(i)
				cursor.Hidden = s.flashHidePeriod
			}

			if s.attributes[i].flash {
				// flash the text
				txt := s.getTextCanvas(i)
				txt.Hidden = s.flashHidePeriod || (s.attributes[i].conceal && !s.revealMode)

				mosaic := s.getMosaicContainer(i)
				mosaic.Hidden = s.attributes[i].conceal && !s.revealMode
				mosaic.Refresh()

			}
		}
		// refreshing the individual elements didn't seem to work properly
		// so refresh the higher level canvas object instead
		s.rootContainer.Objects[0].Refresh()
	}
}

func (s *Screen) Initialise(ctx context.Context, wg *sync.WaitGroup, c fyne.Canvas, textSize float64) *fyne.Container {

	// original container layout
	//s.rootContainer = container.NewWithoutLayout()

	// new container using a grid
	s.rootContainer = container.New(customFyne.NewGridLayout(40))
	s.rootCanvas = c
	s.textSize = float32(textSize)

	//s.rootContainer.Resize(fyne.NewSize(constants.COLS*constants.TEXTWIDTH, constants.ROWS*constants.TEXTHEIGHT))

	for row := 0; row < constants.ROWS; row++ {
		for col := 0; col < constants.COLS; col++ {

			maxContainer := container.New(customFyne.NewMaxLayout())

			//create a background and add it to the character cell
			background := new(canvas.Rectangle)
			maxContainer.Objects = append(maxContainer.Objects, background)

			// mosaic container
			mosaic := customFyne.NewBlock()
			maxContainer.Objects = append(maxContainer.Objects, mosaic.Container)

			//create a text element and add it to the character cell
			txt := canvas.NewText(" ", s.colourMap[7])
			txt.TextSize = s.textSize
			maxContainer.Objects = append(maxContainer.Objects, txt)

			// create the cursot and add it to the character cell
			cursor := new(canvas.Rectangle)
			maxContainer.Objects = append(maxContainer.Objects, cursor)

			s.rootContainer.Objects = append(s.rootContainer.Objects, maxContainer)

		}
	}

	// set the shadow space, this keeps track of control codes
	s.attributes = make([]Attributes, constants.ROWS*constants.COLS, constants.ROWS*constants.COLS)

	// default is with revealMode set
	s.revealMode = false

	// this is for screen writes
	// map ascii to viewdata char set
	s.alphaCharMap = make(map[int]int)
	s.alphaCharMap[0x5f] = 0x23
	s.alphaCharMap[0x5c] = 0xbd // half symbol
	s.alphaCharMap[0x7e] = 0xf7
	s.alphaCharMap[0x7d] = 0xbe // three quarter symbol
	s.alphaCharMap[0x7b] = 0xbc // quarter symbol
	s.alphaCharMap[0x7f] = 0xb6
	s.alphaCharMap[0x23] = 0xa3 // pound sign
	// Colour order by value: r,g,y,b,m,c,w
	colourMap[0] = color.NRGBA{R: 0, G: 0, B: 0, A: 255}       // Black
	colourMap[1] = color.NRGBA{R: 255, G: 0, B: 0, A: 255}     // Red
	colourMap[2] = color.NRGBA{R: 0, G: 255, B: 0, A: 255}     // Green
	colourMap[3] = color.NRGBA{R: 255, G: 255, B: 0, A: 255}   // Yellow
	colourMap[4] = color.NRGBA{R: 0, G: 0, B: 255, A: 255}     // Blue
	colourMap[5] = color.NRGBA{R: 255, G: 0, B: 255, A: 255}   // Magenta
	colourMap[6] = color.NRGBA{R: 0, G: 255, B: 255, A: 255}   // Cyam
	colourMap[7] = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White

	// this ensures that all the screen attributes are set correctly
	s.Clear()
	// run the flash routine as a timed go routine
	wg.Add(1)
	go s.runFlashTimer(ctx, wg)

	s.rootContainer.Refresh()
	return s.rootContainer
}

func (s *Screen) debugWrite(char byte, comment string) {

	if s.Debug {

		var c string
		if char >= 0x20 {
			c = string(char)
		} else {
			c = "."
		}

		//		if s.debugTmp != s.cursorRow {
		//			fmt.Print("\r\nRow/Col\r\n")
		//			s.debugTmp = s.cursorRow
		//		} else {
		//		}
		fmt.Printf("Row:%02d Col:%02d Byte:%02x[%s] - %s\r\n", s.cursorRow, s.cursorCol, char, c, comment)
	}
}
