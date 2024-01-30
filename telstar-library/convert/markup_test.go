package convert

import "testing"

func Test_MarkupToRawV(t *testing.T) {

	type Test struct {
		description string
		input       string
		want        string
	}

	var tests = []Test{
		{"Valid Markup", "[Y]H[F]el[S]o", "\x1b\x43H\x1b\x48el\x1b\x49o"},
		{"Full Set Markup", "[R][G][Y][B][M][C][W][F][S][N][D][-][n][r][g][y][b][m][c][w]", "\x1b\x41\x1b\x42\x1b\x43\x1b\x44\x1b\x45\x1b\x46\x1b\x47\x1b\x48\x1b\x49\x1b\x4c\x1b\x4d\x1b\x5c\x1b\x5d\x1b\x51\x1b\x52\x1b\x53\x1b\x54\x1b\x55\x1b\x56\x1b\x57"},
		{"No Markup", "[Hel[]o", "[Hel[]o"},
		{"New Background", "[R][n][W]hello[C]world[-]Foo Bar",  "\x1b\x41\x1b\x5d\x1b\x47\x68\x65\x6c\x6c\x6f\x1b\x46\x77\x6f\x72\x6c\x64\x1b\x5c\x46\x6f\x6f\x20\x42\x61\x72"},
		{"No Text", "", ""},
	}
	for _, test := range tests {
		if got, err := MarkupToRawV(test.input); got != test.want || err != nil {
			t.Errorf(TEST_ERROR_MESSAGE, test.description)
			bytes := []byte(got)
			print(bytes)
		}
	}
}
/*
"\x1b\x41\x1b\x5d\x1b\x47hello\x1b\x46\x77\x6f\x72\x6c\x64\x1b\x5c\x46\x6f\x6f\x20\x42\x61\x72"
red
new background
white
"hello"


 */