package display

import "testing"

const (
	TEST_ERROR_MESSAGE = "Test Description: \"%s.\""
)

// Double height mosaic test frame
// https://edit.tf/#0:QIECBAgZIEDNAgQIGSBAzQIECBsgQN0CBAgbIEDdAgQIECBAgQMBpNAgQMBgdwNJqECBAMDsBpPAgQcBgdwNJ6ECDwgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQMRpNCgQMRgdyNJqUCBwMDsRpPCgQcRgdyNJ6UCDygQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQMhpNEgQMhgeCNJokCByMDshpPEgQchgeCNJ6kCD0gQIGaAmgQE0CBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQMxpNGgQMxgeENJq0CB0MDsxpNGgQcxgeENJ60CD2gQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQNBpNIgQNBgeGNJrECB2MDtBpPIgQdBgeGNJ7ECD4gQIG40mgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQNRpNCgQNRgeINJrUCB4MDtRpPKgQdRgeINJ7UCD6gQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQNhpNMgQNhgeKNJrkCB6MDthpPMgQdhgeKNJ7kCD8gQIHKAmgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQNxpNOgQNxgeMNJr0CB8MDtxpPOgQdxgeMNJ70CD-gQIECAmgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA

func Test_isLastCol(t *testing.T) {

	type Test struct {
		description string
		inputCol    int
		want        bool
	}

	s := Screen{}

	var tests = []Test{
		{"Col 0", 0, false},
		{"Col 3", 3, false},
		{"Col 22", 22, false},
		{"Col 23", 23, false},
		{"Col 39", 39, true},
		{"Col 40", 40, false},
	}

	for _, test := range tests {
		if got := s.isLastColumn(test.inputCol); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_removeMosaicOffset(t *testing.T) {

	type Test struct {
		description string
		input       byte
		want   byte
	}

	s := Screen{}

	var tests = []Test{
		{"20", 0x20, 0x00},
		{"60", 0x60, 0x20},
	}

	for _, test := range tests {
		if got := s.removeMosaicOffset(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_addMosaicOffset(t *testing.T) {

	type Test struct {
		description string
		input       byte
		want   byte
	}

	s := Screen{}

	var tests = []Test{
		{"00", 0x00, 0x20},
		{"40", 0x20, 0x60},
	}

	for _, test := range tests {
		if got := s.addMosaicOffset(test.input); got != test.want {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_getDoubleHeightMosaic(t *testing.T) {

	type Test struct {
		description string
		input       byte
		wantUpper   byte
		wantLower   byte
	}

	s := Screen{}

	var tests = []Test{
		{"20", 0x20, 0x20, 0x20},
		{"21", 0x21, 0x25, 0x20},
		{"22", 0x22, 0x2A, 0x20},
		{"23", 0x23, 0x2F, 0x20},
		{"24", 0x24, 0x30, 0x21},
		{"25", 0x25, 0x35, 0x21},
		{"26", 0x26, 0x3A, 0x21},
		//{"27", 0x27, 0x30, 0x21},
		//{"28", 0x28, 0x30, 0x21},
		//{"29", 0x29, 0x30, 0x21},
		//{"2A", 0x2a, 0x30, 0x21},
		//{"2B", 0x2b, 0x30, 0x21},
		//{"2C", 0x2c, 0x30, 0x21},
		//{"2D", 0x2d, 0x30, 0x21},
		//{"2E", 0x2e, 0x30, 0x21},
		//{"2F", 0x2f, 0x30, 0x21},
		{"30", 0x30, 0x20, 0x34},
		{"31", 0x31, 0x25, 0x34},
		//{"32", 0x32, 0x70, 0x37},
		//{"33", 0x33, 0x70, 0x37},
		//{"34", 0x34, 0x70, 0x37},
		//{"35", 0x35, 0x70, 0x37},
		//{"36", 0x36, 0x70, 0x37},
		//{"37", 0x37, 0x70, 0x37},
		//{"38", 0x38, 0x70, 0x37},
		//{"39", 0x39, 0x70, 0x37},
		//{"3A", 0x3a, 0x70, 0x37},
		//{"3B", 0x3b, 0x70, 0x37},
		//{"3C", 0x3c, 0x70, 0x37},
		//{"3D", 0x3d, 0x70, 0x37},
		//{"3E", 0x3e, 0x70, 0x37},
		//{"3F", 0x3f, 0x70, 0x37},
		{"60", 0x60, 0x20, 0x68},
		{"61", 0x61, 0x25, 0x68},
		//{"62", 0x62, 0x70, 0x37},
		//{"63", 0x63, 0x70, 0x37},
		//{"64", 0x64, 0x70, 0x37},
		//{"65", 0x65, 0x70, 0x37},
		//{"66", 0x66, 0x70, 0x37},
		//{"67", 0x67, 0x70, 0x37},
		//{"68", 0x68, 0x70, 0x37},
		//{"69", 0x69, 0x70, 0x37},
		//{"6A", 0x6a, 0x70, 0x37},
		//{"6B", 0x6b, 0x70, 0x37},
		//{"6C", 0x6c, 0x70, 0x37},
		//{"6D", 0x6d, 0x70, 0x37},
		//{"6E", 0x6e, 0x70, 0x37},
		{"6F", 0x6f, 0x7f, 0x6b},
		{"70", 0x70, 0x20, 0x7c},
		{"71", 0x71, 0x25, 0x7c},
		//{"72", 0x72, 0x70, 0x37},
		//{"73", 0x73, 0x70, 0x37},
		//{"74", 0x74, 0x70, 0x37},
		//{"75", 0x75, 0x70, 0x37},
		//{"76", 0x76, 0x70, 0x37},
		//{"77", 0x77, 0x70, 0x37},
		//{"78", 0x78, 0x70, 0x37},
		//{"79", 0x79, 0x70, 0x37},
		//{"7A", 0x7a, 0x70, 0x37},
		//{"7B", 0x7b, 0x70, 0x37},
		//{"7C", 0x7c, 0x70, 0x37},
		//{"7D", 0x7d, 0x70, 0x37},
		//{"7E", 0x7e, 0x70, 0x37},
		{"7F", 0x7f, 0x7f, 0x7f},

	}

	for _, test := range tests {
		if gotU, gotL := s.getDoubleHeightMosaic(test.input); gotU != test.wantUpper || gotL != test.wantLower{
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

