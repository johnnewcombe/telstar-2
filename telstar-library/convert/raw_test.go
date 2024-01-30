package convert

import (
	"testing"
)

const (
	BLANK_RAWT = "                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        " +
		"                                        "
)

func Test_RawVToRawT(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        string
		wantLen     int
	}

	// TODO Add tests for BS, HOME, CLS etc
	var tests = []Test{
		{"Page with one line with ESC and CRLF", "\x1b\x531234567890\r\n", ("\x131234567890" + BLANK_RAWT)[:960], 960},
		{"Page with multiple lines", "1234567890\r\nabcdefghijklmnopqrstuvwxyz\r\n", ("1234567890                              abcdefghijklmnopqrstuvwxyz" + BLANK_RAWT)[:960], 960},
		{"Page with multiple lines with missing CR/LF", "1234567890\r\nabcdefghijklmnopqrstuvwxyz", ("1234567890                              abcdefghijklmnopqrstuvwxyz" + BLANK_RAWT)[:960], 960},
		{"Page with one line with CRLF", "1234567890\r\n", ("1234567890" + BLANK_RAWT)[:960], 960},
		{"Page with one line no CRLF", "1234567890", ("1234567890" + BLANK_RAWT)[:960], 960},
	}

	for _, test := range tests {
		if got, err := RawVToRawT(test.input); got != test.want || err != nil || len(got) != test.wantLen {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_createBuffer(t *testing.T) { t.Error("Test not implemented!") }

func Test_RawTToRawV(t *testing.T) { t.Error("Test not implemented!") }

func Test_increaseCol(t *testing.T) {

	type Test struct {
		description string
		inputCol    int
		inputRow    int
		wantCol     int
		wantRow     int
	}
	var tests = []Test{
		{"Col 0, Row 0", 0, 0, 1, 0},
		{"Col 23, Row 0", 23, 0, 24, 0},
		{"Col 38, Row 0", 38, 0, 39, 0},
		{"Col 39, Row 0", 39, 0, 0, 1},
		{"Col 39, Row 23", 39, 23, 0, 0},
	}
	for _, test := range tests {
		if gotCol, gotRow := increaseCol(test.inputCol, test.inputRow); gotCol != test.wantCol || gotRow != test.wantRow {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}
func Test_decreaseCol(t *testing.T) {

	type Test struct {
		description string
		inputCol    int
		inputRow    int
		wantCol     int
		wantRow     int
	}
	var tests = []Test{
		{"Col 0, Row 1", 0, 1, 39, 0},
		{"Col 0, Row 0", 0, 0, 39, 23},
		{"Col 23, Row 0", 23, 0, 22, 0},
	}
	for _, test := range tests {
		if gotCol, gotRow := decreaseCol(test.inputCol, test.inputRow); gotCol != test.wantCol || gotRow != test.wantRow {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}

}

func Test_BlobMerge(t *testing.T) {
	var (
		blob1 string
		blob2 string
		err   error
	)

	type Test struct {
		description string
		inputBlob   string
		inputCol    int
		inputRow    int
		wantBlob    string
		wantBlobLen int
	}
	var tests = []Test{
		{"TestPage", "../tmp/blobs/logo.blob", 0, 0, "", 960},
	}
	for _, test := range tests {

		if blob1, err = EdittfToRawT("https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAxYMcKAbNw6dyCTuyZfCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgKJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRAxQR4s6LSgzEEmdUi0otOogkzo0-lNg1JM-cCg8unNYgQIASBBj67OnXllWIMu7pl5dMOndty7uiDDuyINu_llXLlyBAgQMkE6LXpoIM6IgrxYNSRFpAp2Hpp37sOxBh3ZECBAgQIECAEgQad3TLy3Yemnfuw7EG7L35oMO7Ig75cPTRl5LkCBAgQIEDNBCq05M6LTpoJM6NPpTYNSTPnAoe_bww7vKDDuyIECBAgBIEGLrz07svPmg3Ze_NcDOhI1SnFQTMPTLz6IKHLTjy80CBA0QR59aLSnTYs6ogkzo0-lNg1JM-cgQIECBAgQIECBAgQIEDVBUpQY0aTDQUotCfSqUwU3Dq38kHLfhyIOWXhv5dOaBAgBIEGblv2oJGnPo74fPNBF3Z9mHdkXIECBAgQIECBAgQIECBA2QR4NSLXg2UFOLSrSYcWmCnWKkWYsQVIsWNBsLEEOHMaIASBBh3ZEHTRlQQ9-zfz54diCHh7ZUGHJ2y7unXllXIECBAgQN0EifMkxINmmCqcsPbLsQYd2RBI37NOTD5QbsvfmuQIECBA4CHQc2TDpT50WogcMGCAScgoOnLTi69MqDpvQdNGnmgQIAaBBzy8u2nHlQd9PTQIqZdmXnvzdO-HllE4d2RBt38sq5AgQOQ0yfHQT40ZYgp2adSLNQSZ0aegTIKkWnUQUIMeLTQIECAokSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEFeLBqSItJBGn0osODTqIBIM6EQIEFDDnyoFTJywvoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"); err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
		if blob2, err = EdittfToRawT("https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAv8-dEnz58-fPiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIC7X5_wMlvz581IECBAgQIECBAgQIECBAgQIECBAgQIECBAgLtf6VP47av__UgQIECBAgQIECBAgQIECBAgQIECBAgQIECAu1_tOmjp25u_6BAgQIECBAgQIECBAgQIECBAgQIECBAgQIC7X__df_7X583oECBAgQIECBAgQIECBAgQIECBAgQIECBAgLtf_9w734f__UgQIECBAgQIECBAgQIECBAgQIECBAgQIECAv158__Zam_8_SBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"); err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
		if got, err := RawTMerge(blob1, blob2); err != nil || len(got) != test.wantBlobLen {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
		}
	}
}

func Test_increaseRow(t *testing.T) { t.Error("Test not implemented!") }
func Test_decreaseRow(t *testing.T) { t.Error("Test not implemented!") }
