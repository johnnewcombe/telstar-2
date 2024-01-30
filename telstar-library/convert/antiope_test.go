package convert

import (
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"testing"
)

func Test_RawVToAntiope(t *testing.T) {
	type Test struct {
		description string
		inputBytes  []byte
		wantBytes   []byte
	}
	var tests = []Test{
		//{"Alpha Red Foreground", []byte{0x1b, 0x41, 'H', 'i'}, []byte{0x20, 0x1b, 0x41, globals.SI, 'H', 'i'}},
		{"Alpha Red Foreground, Yellow Background", []byte{0x1b, 0x43, 0x1b, 0x5d, 0x1b, 0x41, 'H', 'i'}, []byte{0x20, 0x1b, 0x43, globals.SI, 0x1b, 0x53, globals.SPC, globals.SI, globals.SPC, 0x1b, 0x41, globals.SI, 'H', 'i'}},
	}
	for _, test := range tests {
		if gotBytes, err := RawVToAntiope(test.inputBytes); !compareBytes(gotBytes, test.wantBytes) ||
			err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_getColumn(t *testing.T) {

	type Test struct {
		description   string
		inputChar     byte
		inputPreChar  byte
		inputCol      int
		wantCol       int
		wantRowChange bool
	}

	var tests = []Test{
		// note that VT and LF on their own dont cause the rowChange flag to be set (it's a compromise!)
		{"Test 1", globals.LF, globals.CR, 4, 4, true},
		{"Test 2", globals.LF, 0x20, 39, 39, false},
		{"Test 3", globals.LF, 0x20, 0, 0, false},
		{"Test 4", globals.DEL, 0x20, 0, 39, true},
		{"Test 5", globals.CR, 0x20, 0, 0, false},
		{"Test 6", globals.CR, 0x20, 39, 0, false},
		{"Test 7", globals.HT, 0x20, 39, 0, true},
		{"Test 8", globals.HT, 0x20, 38, 39, false},
		{"Test 9", globals.VT, 0x20, 39, 39, false},
		{"Test 10", globals.VT, 0x20, 3, 3, false},
		{"Test 11", globals.VT, 0x20, 0, 0, false},
		{"Test 12", globals.BS, 0x20, 0, 39, true},
		{"Test 13", globals.ESC, 0x20, 39, 39, false},
	}
	for _, test := range tests {
		if gotCol, gotRowChange := getColumn(test.inputChar, test.inputPreChar, test.inputCol); gotCol != test.wantCol ||
			gotRowChange != test.wantRowChange {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_processC1Controls(t *testing.T) {

	type Test struct {
		description        string
		inputChar          byte
		previousChar       byte
		inputCurrentColour byte
		wantCurrentColour  byte
		wantBytes          []byte
	}

	var tests = []Test{
		{"Alpha RED to Alpha WHITE", 0x47, 0x1b, 0x41, 0x47, []byte{0x20, globals.ESC, 0x47, globals.SI}},     // G0
		{"Alpha RED to Mosaic WHITE", 0x57, 0x1b, 0x41, 0x57, []byte{0x20, globals.ESC, 0x47, globals.SO}},    // G1
		{"Mosaic WHITE to Mosaic WHITE", 0x57, 0x1b, 0x57, 0x57, []byte{0x20, globals.ESC, 0x57, globals.SO}}, // G1
		{"New Alpha Background.", 0x5d, 0x1b, 0x47, 0x47, []byte{0x20, globals.ESC, 0x57, globals.SI}},        // G0
		{"New Mosaic Background.", 0x5d, 0x1b, 0x57, 0x57, []byte{0x20, globals.ESC, 0x57, globals.SI}},       // G0
	}

	for _, test := range tests {
		if gotBytes, gotCurrentColour, err := processC1Controls(test.inputChar, test.previousChar, test.inputCurrentColour); !compareBytes(gotBytes, test.wantBytes) ||
			gotCurrentColour != test.wantCurrentColour || err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func compareBytes(a []byte, b []byte) bool {

	if a == nil || b == nil || len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != a[i] {
			return false
		}
	}
	return true
}
