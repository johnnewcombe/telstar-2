package convert

import (
	"errors"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/utils"
)

/*
Best Guess so far!

ESC

40 ABK, Alpha Black. Switch to alphabetic, black foreground.
41 ANR, Alpha Red. Switch to alphabetic, red foreground.
42 ANG, Alpha Green. Switch to alphabetic, green foreground.
43 ANY, Alpha Yellow. Switch to alphabetic, yellow foreground.
44 ANB, Alpha Blue. Switch to alphabetic, blue foreground.
45 ANM, Alpha Magenta. Switch to alphabetic, magenta foreground.
46 ANC, Alpha Cyan. Switch to alphabetic, cyan foreground.
47 ANW, Alpha White. Switch to alphabetic, white foreground.
48 FSH, Flashing. Characters displayed flashing between foreground and background.
49 STD, Steady. Terminates flashing.
4a EBX, End Box. Terminates SBX.
4b SBX, Start Box. Defines a non-alphanumeric area, with transparent background. Terminated by EBX.
4c NSZ, Normal Size. Characters normal width and height.
4d DBH, Double Height. Characters normal width and double normal height. Inactive on bottom line.
4e DBW, Double Width. Characters normal height and double normal width. Inactive in last position of line.
4f DBS, Double Size. Characters double height and double normal width. Inactive on bottom line or in last position of line.  (see note 1)
50 MBK, Mosaic Black. Switch to mosaic, black foreground.
51 MSR, Mosaic Red. Switch to mosaic, red foreground.
52 MSG, Mosaic Green. Switch to mosaic, green foreground.
53 MSY, Mosaic Yellow. Switch to mosaic, yellow foreground.
54 MSB, Mosaic Blue. Switch to mosaic, blue foreground.
55 MSM, Mosaic Magenta. Switch to mosaic, magenta foreground.
56 MSC, Mosaic Cyan. Switch to mosaic, cyan foreground.
57 MSW, Mosaic White. Switch to mosaic, white foreground.
58 CDY, Conceal Display. Display characters as spaces (might be terminated by SCD).
59 SPL, Stop Lining. Terminates underlining. For mosaic characters, non-underlined font corresponds to contiguous display, with the blocks within a mosaic character joined together.
5a STL, Start Lining. Begins underlined letters. For mosaics, this corresponds to separated display, with the blocks within a mosaic character shown separated.
5b CSI, Control Sequence Introducer. (see note 2)
5c BBD, Black Background.
5d NBD, New Background. Set background colour to previous foreground colour. The current foreground colour is not affected.
5e HMS, Hold Mosaic. Image subsequently stored control functions as the last received mosaic character.
5f RMS, Release Mosaic. Terminate HMS.

Notes:
[1] I think ESC 4f is double height and width...
[2] CSI: This may be related to ANSI Ctrl sequences.

Example Minitel Page

0C                      // cls
1F 4C 4C                // goto row 12,12
1B 4F                   // double height + double width
1B 43                   // foreground yellow.
46 4F 53 44 45 4D 27    // "FOSDEM'"
1B 4E                   // double width + normal height
0B                      // VT
32 30                   // "20"

TODO Double height extends upwards with Minitel (CEPT p12) but the top row seems to have the attributes.
     May just need to preceed with LF and follow with VT

TODO there may have been a 449 char limit due to how the attributes are stored???
     see file:///Users/john/Development/Repositories/telstar-2/doc/Minitel/reviving-minitel-notes.pdf
     this could mean that many pages transformed from Viewdata are not valid.
     Perhaps the box function in Antiope could get around this?

Notes:
Grahic colour switches to G0 but the escape code should still be sent preceded by an escape

*/

func RawVToAntiope(rawV []byte) ([]byte, error) {

	var (
		//escape       bool
		char          uint8
		previousChar  uint8
		result        []byte
		err           error
		byteSequence  []byte
		currentColour byte
		nextCol       int
		rowChange     bool
	)

	//logger.LogInfo.Printf("Processing %d bytes of RawV data.", len(rawV))
	currentColour = 0x47

	for i := 0; i < len(rawV); i++ {

		char = rawV[i]

		// convert return back to hash (0x23) note that this is not the same as
		// the translations done on characters arriving from the user.
		if char == 0x5f {
			char = 0x23
		}

		// ******************** DEBUG ********************
		//var c             uint8
		//if char < 0x20 {
		//	c = 0x2e
		//} else {
		//	c = char
		//}
		//logger.LogInfo.Printf("%d - Processing byte: 0x%x(%s), Next Column: %d, Row Change: %t ", i, char, string(c), nextCol, rowChange)
		// ******************** DEBUG ********************

		//if previous was esc or New Background
		if previousChar == globals.ESC {

			// we have had an escape (although it hasn't been added to the results yet
			// so we need to process the C1 control char and convert to parallel attributes
			// returns the byteSequence as a []byte
			if byteSequence, currentColour, err = processC1Controls(char, previousChar, currentColour); err != nil {
				logger.LogError.Print(err)
			}

			result = append(result, byteSequence...)

		} else {

			// set the escape flag if appropriate, escape is still sent to the
			// client BUT NOT HERE!, see func processC1Controls()
			if char == globals.ESC {
				previousChar = globals.ESC
				continue
			}
			// TODO Consider the following.
			// If the 'delimiter' is present in a row, it will contain information to change the attributes.
			// Any delimiter character will appear on the screen appearing as an area with background
			// set as the current background color (as a space?).

			result = append(result, char)

		}

		// store the previous char
		previousChar = char

		// set the nextCol and whether a row change should occur,
		nextCol, rowChange = getColumn(char, previousChar, nextCol)
		if rowChange {
			// reset attributes to default
			result = append(result, resetAttributes()...)
			currentColour = 0x47
		}
	}

	return result, nil
}

func AntiopeInputTranslation(inputByte byte) byte {

	// TODO: Look at telstar-client to see translations.
	switch inputByte {
	case 0x23:
		inputByte = 0x5f
	}
	return inputByte
}

func getColumn(char byte, previousChar byte, col int) (int, bool) {

	if char == globals.HOME || char == globals.CLS { // home, cls

		return 0, true

	} else if char == globals.CR {

		return 0, false

	} else if char == globals.LF && previousChar == globals.CR { // simple row change
		// this is a compromise as we are not processing VT and LF as we do not know what the
		// attributes would be at the new location
		return col, true

		//	} else if char == globals.LF || char == globals.VT { // simple row change
		// this is a compromise as we do not know what the attributes should be at the new location
		//		 return col, true

	} else if (char >= globals.SPC && char < globals.DEL) || char == globals.HT { // forward

		// keep track of the output in case of a row change
		col++
		if col > 39 {
			return 0, true
		} else {
			return col, false
		}

	} else if char == globals.BS || char == globals.DEL { // backward

		col--
		if col < 0 {
			return 39, true
		} else {
			return col, false
		}

	} else {
		return col, false // none of the above, no change
	}
}

// processC1Controls returns the Antiope sequence along with a boolean value that indicates
// whether set G1 char set.
func processC1Controls(char byte, previousChar byte, currentColour byte) ([]byte, byte, error) {

	var (
		result []byte
		shift  byte
	)

	// Antiope doesn't use a screen position for C1 controls
	// so we send a space before sending the control. The space will
	// then take on previous attributes.
	if !utils.IsControlC1(char) {
		return result, currentColour, errors.New("The character must be in the fange 0x40-0x5F")
	}

	// FIXME, this is wrong! The colour comes before the New Background.
	//  we need to know what the current alpha or mosaic colour is, it may not be the
	//  previous control, it could have been set earlier in the row

	if char == 0x5d { //new background

		//convert to antiope background colour in range 50h - 58h
		if utils.IsAlphaColour(currentColour) {
			char = currentColour + 0x10
			shift = globals.SI //g1
		} else if utils.IsGraphicColour(currentColour) {
			// char is specifying a graphics background colour this is the same
			// for antiope but we will need to shift to graphics char set
			// so change the char set
			// TODO Consider not sending shift if char already equals the current colour
			char = currentColour
			shift = globals.SO //g1
		} else {
			// not a colour so ignore
		}
		// append char, note that space comes after background colour is set.
		result = append(result, globals.ESC, char, globals.SPC, shift)

		return result, currentColour, nil

	} else if previousChar == globals.ESC {

		// if we have an alpha colour and are not in G0 then switch to G1
		if utils.IsAlphaColour(char) {
			currentColour = char
			shift = globals.SI // G0
		} else if utils.IsGraphicColour(char) {
			// if we have a graphic colour and are not in G1 then switch to G1
			// flip to Antiope foreground colour
			currentColour = char
			char = char - 0x10
			shift = globals.SO // G1
		}

		// Double height with Antiope extends upwards so simulate with LF/VT
		// TODO Need to consider resetting this on a row change otherwise we will be on the wrong line
		if char == 0x4d { // double height
			//result = append(result, globals.LF)
		}

		if char == 0x4c { // normal height
			//result = append(result, globals.VT)
		}

		if char == 0x5c { // end background, output Black Background
			result = append(result, globals.ESC, 0x50, globals.SPC)
		}

		result = append(result, globals.SPC, globals.ESC, char, shift)
	}
	return result, currentColour, nil
}

func resetAttributes() (result []byte) {

	result = append(result, globals.SI)
	result = append(result, []byte(globals.MOSAIC_BLACK)...) // this is black background for Antiope
	result = append(result, []byte(globals.ALPHA_WHITE)...)
	result = append(result, []byte(globals.NORMAL_HEIGHT)...)
	result = append(result, []byte(globals.STEADY)...)
	result = append(result, []byte(globals.STOP_LINING)...)
	result = append(result, []byte(globals.RELEASE_MOSAIC)...)

	return result
}
