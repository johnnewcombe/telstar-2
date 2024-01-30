package convert

import (
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"errors"
	"fmt"
	"strings"
)

/*
Note that there are four different data formats as defined by Telstar.

* RawV - this is 7bit (00-7F) videotex format with control chars between 00-1F and escaped codes for alpha and graphi attributes. Ideally rows should have any trailing spaces removed and replaced with with a CR/LF combination. This will be rendered as is.
* RawT - this is a 24 x 40 block of 7 bits chars (960 chars) in Teletext format (range 00-7F). This format is used internally when manipulating page data. It will be converted to RawV before being rendered.
* Markup - This is Telstar markup and is converted to Raw8 before being rendered.
* EditTf - This is the edit.tf editor's url format. This is converted to RawT before being rendered. Note that as the editTf editor uses 25 rows of 40 columns rather han the standard 24, the last row is ignored.

*/

func RawVToRawT(rawV string) (string, error) {

	var (
		escape bool
		cr     bool
		rawT   strings.Builder
		col    int
		row    int
	)

	pageBuffer := createBuffer(globals.ROWS, globals.COLS)

	for i := 0; i < len(rawV); i++ {

		char := rawV[i]

		if escape {

			// if the previous char was escape and we have a alpha/graphics control char following
			// then un-escape it and set the control char to the 00-1f range as per teletext
			if char >= 0x40 && char <= 0x5f {

				if col == 40 {
					// end of line
					col = 0
					row = increaseRow(row)
				}

				// adjust char to the teletex value and add it to the result
				char -= 0x40
				pageBuffer[row][col] = char
				col, row = increaseCol(col, row)

				// reset the escape flag
				escape = false

			} else if char >= 0x30 && char <= 0x3f {

				// these are control chars that exist in videotex but not in teletex
				// things such as verify, tape control etc. it is assumed that these will
				// not take a printable space so col is not updated

			} else {
				// error, escape should be followed by a char in the range 40-5f
				// some kind of format error
				return "", errors.New("an escape must be followed by a character in the range 40h-5Fh")
			}

		} else if cr {

			// if the previous char was a CR then it must be followed by LF
			if char == globals.LF {

				// reset the char flag
				cr = false

				// move to next row
				col = 0
				row = increaseRow(row)

			} else {
				// error as LF should always follow a CR
				return "", errors.New("a CR must be followed by the character, LF")
			}

		} else if char > 0x1f {

			//normal char
			pageBuffer[row][col] = char
			col, row = increaseCol(col, row)

		} else {

			switch char {
			case globals.NULL:
				// do nothing
			case globals.BS: //BS
				col, row = decreaseCol(col, row)
			case globals.HT: //HT
				col, row = increaseCol(col, row)
			case globals.LF: // LF
				row = increaseRow(col)
			case globals.VT: //VT
				row = decreaseRow(col)
			case globals.CLS:
				//re-initialise the page buffer and set the curreny row/col tp 0/0
				pageBuffer = createBuffer(globals.ROWS, globals.COLS)
				col = 0
				row = 0
			case globals.CR:
				cr = true
			case globals.ESC:
				escape = true
			case globals.HOME:
				col = 0
				row = 0
			}

			/*
				NULL = 0x00
				BS   = 0x08
				HT   = 0x09
				LF   = 0x0A
				VT   = 0x0B
				CLS  = 0x0C
				CR   = 0x0D
				ESC  = 0x1B
				HOME = 0x1E
			*/
		}
	}

	for r := 0; r < 24; r++ {
		for c := 0; c < 40; c++ {

			ch := pageBuffer[r][c]
			if ch == 0 {
				ch = 0x20
			}
			rawT.WriteByte(ch)
		}
	}
	return rawT.String(), nil
}

func createBuffer(rows int, cols int) [][]byte {

	//initialise the slice of slices
	pageBuffer := make([][]byte, rows)
	for i := range pageBuffer {
		pageBuffer[i] = make([]byte, cols)
	}
	return pageBuffer
}

func RawTToRawV(rawT string, rowBegin int, rowEnd int, columnBegin int, columnEnd int, truncateRows bool) (string, error) {

	var (
		result string
		rowNum int
	)

	// TODO implement blobs for titles etc.
	//      perhaps passing in xy and x1y1 params?

	// we are now in teletex mode so extract the portion of the teletex page required
	// col start must be < col end and row begin < row end etc.
	if rowEnd <= rowBegin || columnEnd <= columnBegin {
		errors.New("begin must be less that end for row and columns")
	}
	colsToTake := columnEnd - columnBegin + 1

	// Teletext and Prestel/Telstar is 24 lines, however BBC Mode 7 and edit.tf is 25,
	// In Telstar line 0 is reserved for the header and line 23 (24th line) is reserved
	// for system messages, therefore ignore first line and the last two lines of the raw data.
	for row := rowBegin; row < rowEnd+1; row++ {

		rowNum++

		// TODO: sort out row length... after the following strip
		// get the row but restrict to the row length specified
		rowData := rawT[row*40 : row*40+40][0:colsToTake]

		// if this blob is to be contatenated then the call will probably
		// not want this trimmed
		if truncateRows {
			rowData = strings.TrimRight(rowData, " ")
		}
		rlen := len(rowData)
		for col := columnBegin; col < minOf(colsToTake+1, rlen); col++ {
			asc := rowData[col]

			// for values 00 - 1F, add 40 and precede with an escape
			if asc >= 0x00 && asc <= 0x1f {
				asc += 0x40
				result += "\x1b"
				result += string(asc)
			} else {
				result += string(asc)
			}
		}
		// as rstrip is used for each row, the row could be shorter than 40 chars
		if rlen < 40 {
			result += "\r\n"
		}
	}

	return result, nil

}

func RawTMerge(rawTs ...string) (string, error) {

	// this function takes a set of rawT types (960 chars of Teletext format) and applies them as layers
	// to a rawT background each on top of all the previous. Spaces within the rawT layers are not applied
	// making each layer 'transparent'.

	var (
		result strings.Builder
		index  int
	)
	//initialise the slice of slices
	pageBuffer := createBuffer(globals.ROWS, globals.COLS)
	pageSize := globals.ROWS * globals.COLS

	// main buffer (background) is filled with spaces, in all subsequent blobs (layers) a null character is does not
	// overwrite previous layes i.e. it represents a transparent character
	for r := 0; r < globals.ROWS; r++ {
		for c := 0; c < globals.COLS; c++ {
			pageBuffer[r][c] = 0x20
		}
	}

	for _, rawT := range rawTs {

		size := len(rawT)
		if size != pageSize {
			return "", fmt.Errorf("data item %d (type RawT) is %d bytes in size but should be 960 bytes in size", index, size)
		}

		// copy rawT to pageBuffer for each non-space chars
		for r := 0; r < globals.ROWS; r++ {
			for c := 0; c < globals.COLS; c++ {
				ch := rawT[r*globals.COLS+c]
				if ch != 0x20 {
					pageBuffer[r][c] = rawT[r*globals.COLS+c]
				}
			}
		}

		index++
	}

	// each successive blob is written on top of the previous except where the car is a space
	for r := 0; r < 24; r++ {
		for c := 0; c < 40; c++ {
			ch := pageBuffer[r][c]
			result.WriteByte(ch)
		}
	}
	return result.String(), nil
}

// FIXME Fix this and use instead of MergeRawV in Render module
func RawVMerge(rawVData string, mergeData []string, rows int) (string, error) {

	var (
		err      error
		data     string
		layers   []string
		rawTData string
	)

	// apply any merge-data
	if mergeData != nil {

		// convert main data field to RawT
		if rawTData, err = RawVToRawT(rawVData); err != nil {
			return "", err
		}
		// do the same for each of the merge fields
		for i := 0; i < len(mergeData); i++ {

			// convert all the merge layers to rawT via rawV
			if mergeData[i], err = MarkupToRawV(mergeData[i]); err != nil {
				return "", err
			}
			if mergeData[i], err = RawVToRawT(mergeData[i]); err != nil {
				return "", err
			}
		}

		// combine all of the layers
		layers = append(layers, rawTData)
		layers = append(layers, mergeData...)

		if data, err = RawTMerge(layers...); err != nil {
			return "", err
		}
		// convert main data field back to RawV
		if data, err = RawTToRawV(data, 0, rows-1, 0, 39, true); err != nil {
			return "", err
		}
	}
	return data, nil
}

func increaseCol(col, row int) (int, int) {
	col++
	if col == 40 {
		col = 0
		row = increaseRow(row)
	}
	return col, row
}

func decreaseCol(col, row int) (int, int) {
	col--
	if col < 0 {
		col = 39
		row = decreaseRow(row)
	}
	return col, row
}

func increaseRow(row int) int {
	row++
	if row >= 24 {
		row = 0
	}
	return row
}

func decreaseRow(row int) int {
	row--
	if row < 0 {
		row = 23
	}
	return row
}
