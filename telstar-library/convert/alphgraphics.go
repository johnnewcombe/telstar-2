package convert

import (
	"strings"
)

// for examples see system pages in install.go
const (
	//Original
	//ALPHAGRAPHICS_PAGE = "https://edit.tf/#0:QIECBAgQIIOzhowoI_LDw0acfNBG37uiBAgQIECBAgQIECAugQIGiBAgQIECBogQIHiBAgQIGiBAkQIECRA0QIGiBAgQIC_Pm1bo2rdGlbo2vfm1eoEDdG1bo2rVAgQNUHX0waoEDdu1LrViVasSrViVasSrViRKgQLVjVKgSpUCBA1QJUCVKgQJUqUugQIECBAgQIECBAgQIECBAjRoUCBAgQIEaFAgQIECBAgQIC6BAgQIECBAgQIECBAgQIEDRAgQIECBAgQIECBAgQIECBAgLt0bVujat0bVujat0CDvzYt0CBqgatUDVq1a6_yFqga8-bUulQJVqxK9WJVqxqlQIFixKtQIFqxKjTIVq1KuVpFqxqtWJC6BAgQIECFAgQIEKBAgQIECBAgQIECBAgQIECBAgRo0KBAgLvFjR5sQPNiB4sYPEiB4kQPFjRogaNECBA0QNNCBogQPHjQv14NXq5o1QIGqBq9SIHqRA1wMHqxq1QIEDVA9XNGqBA1atS6VAlWrEq1YlWrEK1IgSoEC1YlSoEqVAgWJUCVAlWpECVKlLoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAu8WNHixo8WNHixo82IHmxAseJGiBo0QNGjRo0QNGiBosWNC7VA1aoGvXg1aoGr1c0WrGiBqgaoGrVA1atWvv-xWvErxYlLpUCVasSpUCBauQpUCVasSoEqBasSo0yFatSpUCVAlQLViQugQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIC-xAgWLGixY0aIEDxYkaIECxY0eLGjxY0eLGiBAgQIECBAgL6kCB4sSrFjVroQLVjR6saIEDV6sarVjVqgaoECBAgQIECAuqQIFqxIsWJVu5IsWJVqxKgQJVqxKgQJVqxKgQIECBAgQIC6BAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgLoECBKgQIECBAgQMECBggQIECBAgQIECBAgQIECBAgQIECAugQIECBAgQIECBAwQIGCBAgQIECBAgQIECBAgQIECBAgQIC6BAgQIECRAgaIECBAgQoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	// Narrowed 1
	//ALPHAGRAPHICS_PAGE = "https://edit.tf/#0:QIECBAgQIIOzhowoI_LDw0acfNBG37uiBAgQIECBAgQIECAugQIGiBAgQIEGhAgQIHiBAgQIGiBAkQIECRA0QIGiBAgQIC_Pm1b60DdWgb60Hf2geoEDfWgb60DVAgQNUHX0waoEDdu1LrViVauQLVyBauQLViBKgQLdyBKqQJUCBA1QJUCVKgQJUqUugQIECBAgQIECBAgQIECBAjRoECBAgQIEaFAgQIECBAgQIC6BAgQIECBAgQIECBAgQIEDRAgQIECBAgQIECBAgQIECBAgLt9aBvrQN9aBvrQN0CDvzQN0CBrqQNUDVq1a6_yFrqQc_aAulVIFq5A9XIFu5AlQIFi5AtQIFq5AjTIVq1KuVpFu5AtWIC6BAgQIECFAgQIkCBAgQIECBAgQIECBAgQIECBAgRo0CBAgLvFjR5sQPNiB4sYPEiB4kQPFjRogaNECBA0QNNCBogQPHjQv14NXq5o1QIGqBq9SIHqRA1wMHqxq1QIEDVA9XNGqBA1atS6VAlWrEq1YlWrEK1IgSoEC1YlSoEqVAgWJUCVAlWpECVKlLoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAu8WNHixo8WNHixo82IHmxAseJGiBo0QNGjRo0QNGiBosWNC7VA1aoGvXg1aoGr1c0WrGiBqgaoGrVA1atWvv-xWvErxYlLpUCVasSpUCBauQpUCVasSoEqBasSo0yFatSpUCVAlQLViQugQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIC-xAgWLGixY0aIEDxYkaIECxY0eLGjxY0eLGiBAgQIECBAgL6kCB4sSrFjVroQLVjR6saIEDV6sarVjVqgaoECBAgQIECAuqQIFqxIsWJVu5IsWJVqxKgQJVqxKgQJVqxKgQIECBAgQIC6BAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgLoECBKgQIECBAgQMECBggQIECBAgQIECBAgQIECBAgQIECAugQIECBAgQIECBAwQIGCBAgQIECBAgQIECBAgQIECBAgQIC6BAgQIECRAgaIECBAgQoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	// Narowed 2
	ALPHAGRAPHICS_PAGE = "https://edit.tf/#0:QIECBAgQIIOzhowoI_LDw0acfNBG37uiBAgQIECBAgQIECAugQIGiBAgQIEGhAgQIHiBAgQIGiBAkQIECRA0QIGiBAgQIC_P2gb60DdWgb60Hf2geoEDfWgb60DVAgQNUHXqgaoEDdu1LrVyBauQLVyBauQLViBKgQLdyBKqQJUCBA1QJVSBKgQJUqUugQIECBAgQIECBAgQIECBAjRoECBAgQIEaFAgQIECBAgQIC6BAgQIECBAgQIECBAgQIEGhAgQIECBAgQIECBAgQIECBAgLt9aBvrQN9aBvrQN0CDvzQa0KBrqQNUDVq1a6_yFrqQc_aAulVIFq5A9XIFu5AlQIFi5AqSIFq5AjTIVq1KuVpFu5AtWIC6BAgQIECFAgQIkCBAgQIECBAgQIECBAgQIECBAgRo0CBAgLvNiB40QPGiDYsYPEiB4kQPGiBpoQNECBA0QNGiBogQPHjQv19IHu1A1QINSBq9SIHqRA14IHu5A1QIEDVA92oGqBA1atS6VUgWrkC1cgXLEK1IgSoEC1cgSqkCVAgWJUCVUgWpECVKlLoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAu82IHmxA82IHmxA8aIHjRBsSIGjRA0aIGjRo0QNGjRAs2IC7XUga6kHX0ga6kD3agW7EGpAgatUDVqgatWvv-x2pUDxcgLpVSBauQJUCBanQJVSBauQKkCBalQK0KBatSpUCVUgQLViAugQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIC-xAgWLGixY0aIEDxYkaIECxY0eLGjxY0eLGiBAgQIECBAgL6kCB4sSrFjVroQLVjR6saIEDV6sarVjVqgaoECBAgQIECAuqQIFqxIsWJVu5IsWJVqxKgQJVqxKgQJVqxKgQIECBAgQIC6BAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgLoECBKgQIECBAgQMECBggQIECBAgQIECBAgQIECBAgQIECAugQIECBAgQIECBAwQIGCBAgQIECBAgQIECBAgQIECBAgQIC6BAgQIECRAgaIECBAgQoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
)

type coordinate struct {
	col int
	row int
}

type BasebandChar [][]int //four rows of three ints
type BasebandText []BasebandChar
type TextRows []string

var alphagraphics = map[uint8]coordinate{
	'a':  coordinate{1, 1},
	'b':  coordinate{4, 1},
	'c':  coordinate{7, 1},
	'd':  coordinate{10, 1},
	'e':  coordinate{13, 1},
	'f':  coordinate{16, 1},
	'g':  coordinate{19, 1},
	'h':  coordinate{22, 1},
	'i':  coordinate{25, 1},
	'j':  coordinate{28, 1},
	'k':  coordinate{31, 1},
	'l':  coordinate{34, 1},
	'm':  coordinate{37, 1},
	'n':  coordinate{1, 5},
	'o':  coordinate{4, 5},
	'p':  coordinate{7, 5},
	'q':  coordinate{10, 5},
	'r':  coordinate{13, 5},
	's':  coordinate{16, 5},
	't':  coordinate{19, 5},
	'u':  coordinate{22, 5},
	'v':  coordinate{25, 5},
	'w':  coordinate{28, 5},
	'x':  coordinate{31, 5},
	'y':  coordinate{34, 5},
	'z':  coordinate{37, 5},
	'A':  coordinate{1, 9},
	'B':  coordinate{4, 9},
	'C':  coordinate{7, 9},
	'D':  coordinate{10, 9},
	'E':  coordinate{13, 9},
	'F':  coordinate{16, 9},
	'G':  coordinate{19, 9},
	'H':  coordinate{22, 9},
	'I':  coordinate{25, 9},
	'J':  coordinate{28, 9},
	'K':  coordinate{31, 9},
	'L':  coordinate{34, 9},
	'M':  coordinate{37, 9},
	'N':  coordinate{1, 13},
	'O':  coordinate{4, 13},
	'P':  coordinate{7, 13},
	'Q':  coordinate{10, 13},
	'R':  coordinate{13, 13},
	'S':  coordinate{16, 13},
	'T':  coordinate{19, 13},
	'U':  coordinate{22, 13},
	'V':  coordinate{25, 13},
	'W':  coordinate{28, 13},
	'X':  coordinate{31, 13},
	'Y':  coordinate{34, 13},
	'Z':  coordinate{37, 13},
	'1':  coordinate{1, 17},
	'2':  coordinate{4, 17},
	'3':  coordinate{7, 17},
	'4':  coordinate{10, 17},
	'5':  coordinate{13, 17},
	'6':  coordinate{16, 17},
	'7':  coordinate{19, 17},
	'8':  coordinate{22, 17},
	'9':  coordinate{25, 17},
	'0':  coordinate{28, 17},
	' ':  coordinate{1, 21},
	'\'': coordinate{4, 21},
	'.':  coordinate{7, 21},
	',':  coordinate{10, 21},
	':':  coordinate{13, 21},
	';':  coordinate{16, 21},
}

func TextToAlphagraphics(text string, graphicColour string, htab int, vtab int) (string, error) {

	// This returns a blob of videotext representing the text in alphagraphics
	// which, when rendered, will appear at the specific =vtab, htab position

	var (
		textRows TextRows
		err      error
	)
	// if the escape has already been added such as when theme constants are
	// passed in
	if strings.HasPrefix(graphicColour, "\x1b") {
		graphicColour = graphicColour[1:]
	}

	if textRows, err = getTextRows(text); err != nil {
		return "", err
	}

	return rowsToBlob(textRows, graphicColour, htab, vtab)
}

func getTextRows(text string) (TextRows, error) {

	var (
		textData BasebandText
		err      error
	)
	if textData, err = getTextData(text); err != nil {
		return []string{}, err
	}

	// text_data now contains each of the characters as 12 integers representing the character (4 rows of 3)
	// in 'baseband' (teletext) format i.e. the characters are in the range 0x00 0x3F. Before being transmitted, those
	// chars between 0x20 and 0x3F will need to have 0x40 added and those between 0x00 and 0x1F will need
	// to have 0x20 added.

	bitShiftRequired := false
	for i := 0; i < len(text); i++ {

		charData := textData[i]

		if bitShiftRequired {
			if charData, err = shiftPixelsRight(charData); err != nil {
				return []string{}, err
			}
		}
		if charData, err = proportionallySpace(charData); err != nil {
			return []string{}, err
		}

		bitShiftRequired = !hasRightBorder(charData)

	}

	return videotexEncode(textData)

}

func getTextData(text string) (BasebandText, error) {

	// note that maps are reference types so no pointer needed
	var (
		err                   error
		alphagraphicsPageData string
		baseBandText          BasebandText
	)

	//baseBandChar will contains each of the characters as 12 integers representing the character (4 rows of 3)
	// e.g. x =3 and y = 4
	var basebandChar BasebandChar

	if alphagraphicsPageData, err = EdittfToRawT(ALPHAGRAPHICS_PAGE); err != nil {
		return baseBandText, err
	}
	for i := 0; i < len(text); i++ {
		key := text[i]
		coords := alphagraphics[key]

		// get the character data
		if basebandChar, err = getCharData(alphagraphicsPageData, coords); err != nil {
			return baseBandText, err
		}
		// add it to the result
		baseBandText = append(baseBandText, basebandChar)

	}
	return baseBandText, nil
}

func getCharData(alphagraphicsPageData string, coords coordinate) (BasebandChar, error) {

	var (
		basebandChar BasebandChar
		rowNum       int
		colNum       int
		rowData      string
	)

	dy := 4
	dx := 3

	//initialise the slice of slices
	basebandChar = make([][]int, dy)
	for i := range basebandChar {
		basebandChar[i] = make([]int, dx)
	}

	row1 := coords.row + 4 // row start/end
	col1 := coords.col + 3 // col start end

	// Useful for HEX dumping data
	//fmt.Printf("Teletex:\n%s\n", hex.Dump([]byte(alphagraphicsPageData)))

	for r := coords.row; r < row1; r++ {

		rowData = alphagraphicsPageData[r*40 : r*40+40]
		colNum = 0

		for c := coords.col; c < col1; c++ {

			//this gives us the raw data byte but...
			i := rowData[c]

			if i >= 0x60 && i <= 0x7f {
				i -= 0x40
			} else if i >= 0x20 && i <= 0x3f {
				i -= 0x20
			} else {
				i = 0
			}

			basebandChar[rowNum][colNum] = int(i)
			colNum++
		}
		rowNum++
	}

	return basebandChar, nil
}

func shiftPixelsRight(charData BasebandChar) (BasebandChar, error) {

	// TODO: if char without border doesn't use row 3, total_cols-1 for decenders then we only need to bitshift any chars that
	//       use bits 0, 2 and 4 in rows 0, 1 and 2. e.g. a 'j' can follow an f without the need to bit shift
	cols := len(charData[0])

	for row := 0; row < 4; row++ {

		// loop through the columns in reverse order
		for i := cols - 1; i >= 0; i-- {
			if i < cols-1 {
				// save the val of bits 1, 3 and 5
				tmp := charData[row][i] & 0x2a

				// shift left to make them bits 0, 2 and 4
				tmp = tmp >> 1
				charData[row][i+1] = charData[row][i+1] | tmp

			}
			// set to only the value of bits 0, 2 and 4
			charData[row][i] = charData[row][i] & 0x15

			// shift left to move the valuesto bits 1, 3 and 5
			charData[row][i] = charData[row][i] << 1

			if i == 0 {
				// all shifted so blank the left hand border (pixels 0, 2 and 4 of col 0) by including only cols 1, 3 and 5
				charData[row][i] = charData[row][i] & 0x2a
			}
		}
	}

	return charData, nil
}

func proportionallySpace(charData BasebandChar) (BasebandChar, error) {

	// assume char is minimum width (always something in col 0)
	cols := 1

	// traverse all but the first column
	for col := 1; col < 3; col++ {

		// assume we can remove this column in order to proportionally space
		blankCol := true

		// check that each row in this column is blank
		for row := 0; row < 4; row++ {

			if charData[row][col] != 0 {
				// row ot blank so this column must stay
				blankCol = false // can"t remove column
				break
			}
		}

		if !blankCol {
			cols += 1
		}
	}

	// we have a blank column so reduce the length of each row
	for i := 0; i < 4; i++ {
		// TODO need to reduce the length of each row to the value of cols
		//  may need to switch to slices rather than fixed array?
		charData[i] = charData[i][:cols]

	}

	// check for blank right hand pixel border. This is used as a space characters, if a border does not exist the following
	// character will need to be shifted over a pixel. Bits of the byte represent the pixels as follows
	//   -------
	//    0 | 1
	//   -------
	//    2 | 3
	//   -------
	//    4 | 5
	//   -------
	//
	// Therefore, to check for a right hand border (i.e. no pixels set) simply check bits 1, 3 and 5 for each row

	return charData, nil
}

func hasRightBorder(charData BasebandChar) bool {
	cols := len(charData[0])
	tmp := 0
	for i := 0; i < 4; i++ {
		tmp += charData[i][cols-1] & 0x2a // bits 1, 3 and 5
	}

	// return true if border exists
	return tmp == 0
}

func videotexEncode(textData BasebandText) ([]string, error) {

	var (
		charData BasebandChar
	)
	// result defined for four rows of text
	result := make([]string, 4)

	for i := 0; i < len(textData); i++ {
		charData = textData[i]
		cols := len(charData[0])

		for r := 0; r < 4; r++ {
			for c := 0; c < cols; c++ {
				char := charData[r][c]

				if char >= 0 && char <= 0x1f {
					result[r] += string(rune(char + 0x20))
				} else if char >= 0x20 && char <= 0x3f {
					result[r] += string(rune(char + 0x40))
				}
			}
		}
	}

	// text dumps for debugging
	//fmt.Printf("textRows 0:\n%s\n", hex.Dump([]byte(result[0])))
	//fmt.Printf("textRows 1:\n%s\n", hex.Dump([]byte(result[1])))
	//fmt.Printf("textRows 2:\n%s\n", hex.Dump([]byte(result[2])))
	//fmt.Printf("textRows 3:\n%s\n", hex.Dump([]byte(result[3])))

	return result, nil
}

func rowsToBlob(textRows TextRows, graphicColour string, htab int, vtab int) (string, error) {

	var (
		blb string
		ht  string
	)

	for v := 0; v < vtab; v++ {
		blb += "\x0a"
	}
	for h := 0; h < htab; h++ {
		ht += "\x09"
	}

	for r := 0; r < len(textRows); r++ {
		blb += ht + "\x1b" + graphicColour + textRows[r] + "\r\n"
	}

	return blb, nil
}
