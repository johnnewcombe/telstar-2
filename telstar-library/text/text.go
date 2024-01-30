package text

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	htmlTagStart    = 60 // Unicode `<`
	htmlTagEnd      = 62 // Unicode `>`
	formatSeparator = "\r\n"
)

// StringCharToRune converts a single char length string to rune if it cannot convert, 0 is returned
func StringCharToRune(char string) rune {
	if len(char) > 0 {
		return []rune(char)[0]
	} else {
		return 0
	}
}

// StringToRune converts a string to a slice of runes.
func StringToRune(s string) []rune {
	return []rune(s)
}

// RuneToString Converts a single rune to a string
func RuneToString(r rune) string {
	return string(r)
}

//RunesToString converts a slice of runes to a string.
func RunesToString(r []rune) string {
	return string(r)
}

//GetMarkupLen returns length of text after any markup has been evaluated.
func GetMarkupLen(markupText string) int {

	// See https://regex101.com/

	/*
			FIXME: Need to account for Home [@], VT [V] LF, CR etc..
		     maybe just raise an error?
			These sequences all need to be counted as one char
			[n],[r],[g],[y],[b],[m],[c],[w],[R],[G],[Y],[B],[M],[C],[W],[F],[S],[N],[D],[-]


			These sequences (horizontal lines) all need to be counted as 39 chars
			[h.],[m.],[l.],[h-],[m-],[l-]

			These sequences (cursor on/off) all need to be counted as ZERO chars
			[_+],[_-]

			The following placeholders could exist also (some are used in the telstar-rss program)
			and should be ignored i.e. counted as ZERO chars
			[NAME],[GREETING],[SERVER],[DATE],[TIME]
	*/

	// address all the markup that will be a single character when decoded and replace with a dot
	m1 := regexp.MustCompile(`(?:\[[nrgybmcwRGYBMCWFSNDH]\]|\[-\])`) //\[[hml].\]|\[[hml]\]-|\[-\])`)
	tmp := m1.ReplaceAllString(markupText, ".")

	// each of these this needs to add 39 to the count
	m2 := regexp.MustCompile(`(?:\[[hml].\]|\[[hml]-\])`)
	tmp = m2.ReplaceAllString(tmp, ".......................................")

	// remove those that do not take up a char or should be ignored
	m3 := regexp.MustCompile(`(?:\[_\+\]|\[_-\]|\[NAME\]|\[GREETING\]|\[SERVER\]|\[DATE\]|\[TIME\])`)
	tmp = m3.ReplaceAllString(tmp, "")

	return len(tmp)

}

// GetDisplayLen returns length of viewdata string, CR and LF are
// ignored but HT and BS are accounted for
func GetDisplayLen(text string) int {

	var count int

	for i := 0; i < len(text); i++ {
		// check for non-control and tab
		if text[i] >= 32 || text[i] == 0x9 {
			count++
		}
		// check for BS
		if text[i] == 0x8 {
			if count > 0 {
				count--
			}
		}
	}

	return count
}

func AreEqualByteSlices(bytes1 []byte, bytes2 []byte) bool {

	// using explicit variable so that breakpoints can be added during testing
	var result bool
	result = bytes.Compare(bytes1, bytes2) == 0
	return result
}

func PadTextLeft(text string, newLength int) string {

	count := GetDisplayLen(text)
	if count > newLength {
		//FIXME concatenate if too big see PadTextRight
		return text
	}

	pad := newLength - count
	if pad > 0 {

		for i := 0; i < pad; i++ {
			text = " " + text
		}
	}
	return text
}

func PadTextRight(text string, newLength int) string {

	count := GetDisplayLen(text)
	if count > newLength {
		return text[:newLength]
	}

	pad := newLength - count
	if pad > 0 {

		for i := 0; i < pad; i++ {
			text += " "
		}
	}
	return text
}

//Format Returns the specified text in rows of given length split by words, and the number of rows.
func Format(text string, cols int) ([]string, int) {
	if cols > 0 {

		fstring := formatString(forceASCII(cleanText(text)), cols)
		result := strings.Split(fstring, "\r\n")
		return result, len(result)
	}
	return []string{}, 0
}

func forceASCII(s string) string {
	rs := make([]rune, 0, len(s))
	for _, r := range s {
		if r <= 127 {
			rs = append(rs, r)
		}
	}
	return string(rs)
}

// cleanText aggressively strips HTML tags from a string.
// It will keep anything between `>` and `<`.
func cleanText(s string) string {

	// Special case (occurs in private eye
	s = strings.ReplaceAll(s, "&lt;p&gt;", "")
	s = strings.ReplaceAll(s, "</p>", "</p> ")
	s = strings.ReplaceAll(s, "</P>", "</P> ")

	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character and we are not in an HTML tag, save it.
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()

	// replace paras
	s = strings.Replace(s, "\r\n\r\n", "|#|", -1)

	// replace cr's
	s = strings.Replace(s, "\r\n ", "\r\n", -1)
	//s = strings.Replace(s, "\n", " ", -1)

	// replace multiple spaces
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	// put paras back
	s = strings.Replace(s, "|#|", "\r\n\r\n", -1)

	// swap common elements
	s = strings.Replace(s, "\r\n ", "\r\n", -1)
	s = strings.Replace(s, " :", ":", -1)
	s = strings.Replace(s, " .", ".", -1)
	s = strings.Replace(s, " ,", ",", -1)
	s = strings.Replace(s, " |", ":", -1)

	//|
	// THESE LOOK THE SAME BUT ARE DIFFERENT UNICODE CHARS!!!
	s = strings.Replace(s, "’", "'", -1)
	s = strings.Replace(s, "‘", "'", -1) // not the same as the one above see below
	s = strings.Replace(s, "‘", "'", -1) // not the same as the one above see below
	s = strings.Replace(s, "“", "'", -1)
	s = strings.Replace(s, "”", "'", -1) // not the same as the one above see below

	/* check the results of each !
	a:=[]byte("’")
	b:=[]byte("‘")
	c:=[]byte("‘")
	d:=[]byte("“")
	e:=[]byte("”")
	*/

	// this often appears in links and is simply removed
	s = strings.Replace(s, "Continue reading...", "", -1)
	return s
}

func formatString(text string, cols int) string {

	var (
		words, pWords []string
		sb            = strings.Builder{}
	)

	//fmt.Println(text)
	words = strings.Split(text, " ")
	for _, w := range words {

		if len(w) > cols-2 {
			subWords := strings.Split(w, "-")
			for _, sw := range subWords {
				//fmt.Printf("%s ", sw)
				if len(sw) > cols-2 {
					pWords = append(pWords, fmt.Sprintf("%s...", sw[:cols-3]))
				}
				pWords = append(pWords, sw)
			}
		}
		pWords = append(pWords, w)
	}

	// format the processed words
	var line string
	for _, w := range pWords {

		l := len(w)

		// if the current word will fit into the row then add it
		if len(line)+l < cols {
			line += fmt.Sprintf("%s ", w)
		} else {
			//not enough space in line so add current line to output and start a new one
			line += formatSeparator
			sb.WriteString(line)
			line = fmt.Sprintf("%s ", w)
		}
	}
	sb.WriteString(line)
	return sb.String()
}

func CleanUtf8(s string) string {

	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s = string(v)
	}
	return s
}
