package convert

import (
	"errors"
	"fmt"
	"strings"
)

// EdittfToRawT converts from an Edit.tf format URL to a RawT format.
func EdittfToRawT(edittfUrl string) (string, error) {

	const ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	var (
		cbit        int
		cpos        int
		encodedData string
		err         error
	)

	if encodedData, err = GetEditTfRawData(edittfUrl); err != nil {
		return "", err
	}

	decodedData := make([]byte, 1000)
	var count = 0

	for index := 0; index < len(encodedData); index++ {
		findex := strings.Index(ALPHABET, string(encodedData[index]))

		if findex == -1 {
			return "", fmt.Errorf("the encoded character at position %d could not be found in the alphabet", findex)
		}
		for b := 0; b < 6; b++ {
			bit := findex & (1 << (5 - b))
			if bit > 0 {
				cbit = (index * 6) + b
				cpos = cbit % 7.0
				cloc := int(float64(cbit-cpos) / 7.0)
				decodedData[cloc] |= 1 << (6 - cpos)
			}
		}
		count++
	}

	return string(decodedData), nil
}

// EdittfToRawV converts from an Edit.tf format URL to a RawV format.
// The Edit.tf data is always truncated to 24 rows in line with the Teletext and Viewdata formats.
// However, the conversion to RawV can be performed on a subset of rows by specifying the
// rowBegin and rowEnd. These values are inclusive, e.g. rows 1 to 22 will return 22 rows of data.
// The function returns a character string along with any error encountered.
func EdittfToRawV(edittfUrl string, rowBegin int, rowEnd int, truncateRows bool) (string, error) {

	var (
		err  error
		rawT string
	)

	if rawT, err = EdittfToRawT(edittfUrl); err != nil {
		return "", err
	}
	return RawTToRawV(rawT[:960], rowBegin, rowEnd, 0, 39, truncateRows)
}

func minOf(vars ...int) int {

	// e.g. MinOf(3, 9, 6, 2) returns 2

	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func maxOf(vars ...int) int {

	// e.g. MaxOf(3, 9, 6, 2) returns 9

	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func GetEditTfRawData(url string) (string, error) {

	data := strings.Split(url, ":")

	if len(data) < 2 {
		return "", errors.New("the url is invalid")
	}

	rawEdata := data[2]
	rawEdata = strings.TrimRight(rawEdata, "\n")
	rawEdata = strings.TrimRight(rawEdata, "\r")

	length := len(rawEdata)
	if !(length == 1120 || length == 1167) {
		return "", fmt.Errorf("encoded frame should be exactly 1120 or 1167 characters in length, actual length was %d", length)
	}

	return rawEdata, nil
}
