package convert

import (
	"testing"
)

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

func Test_textToAlphagraphics(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_getTextRows(t *testing.T) {
	t.Error("Test not implemented!")
}
func Test_getCharData(t *testing.T) {

	var (
		alphagraphicsPageData string
		err                   error
	)

	if alphagraphicsPageData, err = EdittfToRawT(ALPHAGRAPHICS_PAGE); err != nil {
		t.Errorf(TEST_ERROR_MESSAGE, "cannot load alphagraphicsPageData needed for the test")
	}

	type Test struct {
		description string
		input       string
		inputCoord  coordinate
		want        BasebandChar
	}

	var tests = []Test{ // for these tests we are passing a coordinate for data to be extracted from the alphagraphicsPageData
		{"The letter 'H'", alphagraphicsPageData, coordinate{22, 9}, BasebandChar{[]int{20, 40, 0}, []int{29, 46, 0}, []int{5, 10, 0}, []int{0, 0, 0}}},
		{"The letter 'e'", alphagraphicsPageData, coordinate{13, 1}, BasebandChar{[]int{0, 0, 0}, []int{55, 59, 0}, []int{13, 12, 0}, []int{0, 0, 0}}},
		{"The letter 'a'", alphagraphicsPageData, coordinate{1, 1}, BasebandChar{[]int{0, 0, 0}, []int{51, 59, 0}, []int{13, 14, 0}, []int{0, 0, 0}}},
		{"The letter 'd'", alphagraphicsPageData, coordinate{10, 1}, BasebandChar{[]int{0, 40, 0}, []int{23, 43, 0}, []int{13, 14, 0}, []int{0, 0, 0}}},
	}

	for _, test := range tests {
		if got, err := getCharData(test.input, test.inputCoord); !compareBasebandChar(test.want, got) || err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_getTextData(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        BasebandText
	}

	var tests = []Test{
		{"The word 'Head'", "Head",
			BasebandText{
				BasebandChar{[]int{20, 40, 0}, []int{29, 46, 0}, []int{5, 10, 0}, []int{0, 0, 0}},
				BasebandChar{[]int{0, 0, 0}, []int{55, 59, 0}, []int{13, 12, 0}, []int{0, 0, 0}},
				BasebandChar{[]int{0, 0, 0}, []int{51, 59, 0}, []int{13, 14, 0}, []int{0, 0, 0}},
				BasebandChar{[]int{0, 40, 0}, []int{23, 43, 0}, []int{13, 14, 0}, []int{0, 0, 0}},
			},
		},
	}

	for _, test := range tests {
		if got, err := getTextData(test.input); !compareBasebandText(test.want, got) || err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_hasRightBorder(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_videotexEncode(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_proportionallySpace(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_shiftPixelsRight(t *testing.T) {
	t.Error("Test not implemented!")
}

/* this function no longer exists
func Test_textToAlphgraphics(t *testing.T) {

	type Test struct {
		description  string
		input        string
		wantSrc      string
		wantTextRows []string
	}

	var tests = []Test{
		{"Headlines", "Headlines", "Headlines",[]string{"","",""}},
		{"Head", "Head", "Head",[]string{"","",""}},
	}
	for _, test := range tests {
		if got, err := TextToAlphgraphics(test.input); !compareTextRows(got.textRows, test.wantTextRows) ||
			got.sourceText != test.wantSrc ||
			err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
*/
func compareTextRows(textRows1 []string, textRows2 []string) bool {

	if len(textRows1) != len(textRows2) {
		return false
	}

	for row := 0; row < len(textRows1); row++ {
		if textRows1[row] != textRows2[row] {
			return false
		}
	}

	return true
}

func compareBasebandText(text1 BasebandText, text2 BasebandText) bool {

	// check length
	if len(text1) != len(text2) {
		return false
	}

	// check each char
	for char := 0; char < len(text1); char++ {
		if !(compareBasebandChar(text1[char], text2[char])) {
			return false
		}
	}
	return true
}

func compareBasebandChar(char1 BasebandChar, char2 BasebandChar) bool {

	// check length
	if len(char1) != len(char2) {
		return false
	}

	// check ints of the char
	for row := 0; row < len(char1); row++ {
		for col := 0; col < len(char1[row]); col++ {
			if char1[row][col] != char2[row][col] {
				return false
			}
		}
	}
	return true
}
