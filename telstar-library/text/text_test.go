package text

import "testing"

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

// CharToRune converts a single char length string to rune if it cannot convert, 0 is returned
func Test_CharToRune(t *testing.T) {
	t.Error("Test not implemented!")
}

// StringToRune converts a string to a slice of runes.
func Test_StringToRune(t *testing.T) {
	t.Error("Test not implemented!")
}

// RuneToString Converts a single rune to a string
func Test_RuneToString(t *testing.T) {
	t.Error("Test not implemented!")
}

// RunesToString converts a slice of runes to a string.
func Test_RunesToString(t *testing.T) {
	t.Error("Test not implemented!")
}

func Test_GetMarkupLen(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        int
	}

	var tests = []Test{
		{"Test 1", "[r][y][-]", 3},
		{"Test 2", "[Q][y][-]", 5},
		{"Test 3", "abc[][_-][_+]", 5},
		{"Test 4", "abc[][_-_+]", 11},
		{"Test 5", "[r][y][-][SERVER]", 3},
		{"Test 6", "[Q][y][-][GREETING][SERVER]", 5},
		{"Test 7", "[NAME]abc[][DATE][_-][_+]", 5},
		{"Test 8", "[CONTENT]abc[][_-_+][H]", 21},
		{"Test 9", "[l.]", 39},
	}

	// run tests
	for _, test := range tests {
		if got := GetMarkupLen(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_GetDisplayLen(t *testing.T) {

	// TODO: Implement Test
	t.Error("Test not implemented!")

	type Test struct {
		description string
		input       string
		want        int
	}

	var tests = []Test{
		{"Test 1", "", 3},
		{"Test 2", "", 5},
		{"Test 3", "", 5},
		{"Test 4", "", 11},
	}

	// run tests
	for _, test := range tests {
		if got := GetDisplayLen(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_AreEqualByteSlices(t *testing.T) {

	// define tests
	type Test struct {
		description string
		inputSlice1 []byte
		inputSlice2 []byte
		want        bool
	}

	var tests = []Test{
		{"Slices Equal", []byte{9, 10, 1, 2, 3}, []byte{9, 10, 1, 2, 3}, true},
		{"Slices Not Equal", []byte{9, 10, 1, 2, 3, 4}, []byte{9, 10, 1, 2, 3}, false},
	}

	// run tests
	for _, test := range tests {
		if got := AreEqualByteSlices(test.inputSlice1, test.inputSlice2); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_PadTextRight(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_Format(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_RemoveTextBetween(t *testing.T) {

	type Test struct {
		description string
		input       string
		start       string
		end         string
		want        string
	}

	var tests = []Test{
		{"Test 1", "Hello testing World", "Hello", " World", "Hello World"},
		{"Test 2", "Hello test World", "Hello", " World", "Hello World"},
		{"Test 3", "Hello testing world", "Hello", " World", "Hello testing world"},
		{"Test 4", "The post <a href=\"\"/> appeared", "The post", " appeared", "The post appeared"},
	}

	// run tests
	for _, test := range tests {
		if got := RemoveTextBetween(test.input, test.start, test.end); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_cleanText(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        string
	}

	var tests = []Test{
		{"Test 1", "Hello   World", "Hello World"},
		{"Test 2", "Hello   \r\nWorld", "Hello World"},
		{"Test 3", "Hello   \r\n\r\nWorld", "Hello \r\n\r\nWorld"},
		{"Test 4", "Hello   World :", "Hello World:"},
		{"Test 5", "Hello   \tWorld .", "Hello World."},
		{"Test 6", "‘Hello‘ World .", "'Hello' World."},
		{"Test 7", "<p>Hello</p><p>World</p>.", "Hello World."},
	}
	//‘
	// run tests
	for _, test := range tests {
		if got := cleanText(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_formatString(t *testing.T) {
	// TODO: Implement Test
	t.Error("Test not implemented!")
}

func Test_CleanUtf8(t *testing.T) {

	// Test with this
	//s := "a\xc5z"
	// define tests
	type Test struct {
		description string
		input       string
		want        string
	}

	var tests = []Test{
		{"a\xc5z", "a\xc5z", "az"},
	}

	// run tests
	for _, test := range tests {
		if got := CleanUtf8(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}
