package utils

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wagslane/go-password-validator"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	REGEXPAGEID = "^[0-9]+[a-z]$"
	REGEXUSERID = "^[0-9]+$"
)

func PageInScope(basePage int, pageNumber int) bool {

	// base page 0 has access to everything.
	if basePage == 0 {
		return true
	}
	base := strconv.Itoa(basePage)
	page := strconv.Itoa(pageNumber)

	if strings.HasPrefix(page, base) {
		return true
	}
	return false
}

func CreateGuid() string {
	uuidWithHyphen := uuid.New()
	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

func IsValidPageId(pageId string) bool {

	if len(pageId) > 10 {
		return false
	}

	regExFrame, err := regexp.Compile(REGEXPAGEID)
	if err != nil {
		return false
	}
	return regExFrame.MatchString(pageId)
}

func IsValidFrameId(frameId rune) bool {

	// returns true if an ascii value for lowercase a-z
	return frameId >= 0x61 && frameId <= 0x7A

}

func IsValidUserId(userId string) bool {

	if len(userId) != 10 {
		return false
	}

	regExFrame, err := regexp.Compile(REGEXUSERID)
	if err != nil {
		return false
	}
	return regExFrame.MatchString(userId)
}

func IsValidRoutingTable(routingTable []int) bool {

	if len(routingTable) == 11 {

		//check routing table entries
		for _, i := range routingTable {

			s := fmt.Sprintf("%d", i)
			if len(s) > 9 {
				return false
			}
		}
		return true
	}
	return false
}

func IsValidPageNumber(pageNumber int) bool {

	// returns true if an ascii value for lowercase a-z
	return pageNumber >= 0 && pageNumber <= 999999999

}

func CreateDefaultRoutingTable(pageNumber int) []int {

	var routingTable []int

	routingTable = make([]int, 11)
	for i := 0; i < len(routingTable)-1; i++ {
		routingTable[i] = pageNumber*10 + i
	}
	routingTable[10] = pageNumber / 10
	return routingTable
}

func ConvertPidToPageId(pageNumber int, frameId string) (string, error) {

	// check if no frame id
	if len(frameId) == 0 {
		frameId = "a"
	}

	pageId := fmt.Sprintf("%d%s", pageNumber, frameId)

	if IsValidPageId(pageId) {
		return pageId, nil
	} else {
		return "", errors.New("invalid page id")
	}
}

func ConvertPageIdToPID(pageId string) (int, string, error) {

	var (
		frameId    string
		pageNumber int
		//result     PID
		err error
	)

	if IsValidPageId(pageId) {

		//result.pageNumber, err = strconv.Atoi(pageId[0 : len(pageId)-1])
		pageNumber, err = strconv.Atoi(pageId[0 : len(pageId)-1])
		if err != nil {
			return 0, "", fmt.Errorf("unable to convert Page ID: %s.", pageId)
		}
		//result.frameId = pageId[len(pageId)-1:]
		frameId = pageId[len(pageId)-1:]

	} else {
		return 0, "", fmt.Errorf("unable to convert Page ID: %s.", pageId)
	}

	return pageNumber, frameId, nil
}

func ConvertPageIdToPID2(pageId string) (int, rune, error) {

	var (
		frameId    rune
		pageNumber int
		err        error
	)

	if IsValidPageId(pageId) {

		//result.pageNumber, err = strconv.Atoi(pageId[0 : len(pageId)-1])
		pageNumber, err = strconv.Atoi(pageId[0 : len(pageId)-1])
		if err != nil {
			return 0, 0, fmt.Errorf("unable to convert Page ID: %s.", pageId)
		}
		//result.frameId = pageId[len(pageId)-1:]
		frameId = []rune(pageId[len(pageId)-1:])[0]

	} else {
		return 0, 0, fmt.Errorf("unable to convert Page ID: %s.", pageId)
	}

	return pageNumber, frameId, nil
}

func IsNumeric(char byte) bool {

	// returns the ordinal value of the first character of a string
	// returns true if an ascii value for 0 - 9
	//var ord = int([]rune(char)[0])
	return char >= 0x30 && char <= 0x39

}

func IsAlphaNumeric(char byte) bool {

	// returns the ordinal value of the first character of a string
	// returns true if an ascii value for 32 - 127 and BS
	//var ord = int([]rune(char)[0])
	return (char >= 0x20 && char <= 0x7F) || char == 0x08

}

func IsGraphicCharacter(char byte) bool {
	return char >= 0x20 && char <= 0x3f || char >= 0x60 && char <= 0x7f
}

func IsControlC0(char byte) bool {
	return char >= 0x00 && char <= 0x1f
}

func IsControlC1(char byte) bool {
	return char >= 0x40 && char <= 0x5f
}

func IsGraphicColour(char byte) bool {
	return char >= 0x50 && char <= 0x58
}

func IsAlphaColour(char byte) bool {
	return char >= 0x40 && char <= 0x48
}

func GetFollowOnFrameId(frameId rune) (rune, error) {

	var (
		ordinalFrameId byte
	)
	ordinalFrameId = byte(frameId)

	if ordinalFrameId < 122 {
		ordinalFrameId++
	} else {
		ordinalFrameId = 97
	}
	return rune(ordinalFrameId), nil
}

// FIXME This should be a rune
func GetFollowOnPageId(pageId string) (string, error) {
	var (
		pageNumber     int
		frameId        string
		ordinalFrameId byte
		err            error
	)

	if !IsValidPageId(pageId) {
		return "", errors.New("pageId is invalid")
	}
	if pageNumber, frameId, err = ConvertPageIdToPID(pageId); err != nil {
		return "", err
	}
	//ordinalFrameId = int([]rune(frameId)[0])
	ordinalFrameId = []byte(frameId)[0]

	if ordinalFrameId < 122 {
		ordinalFrameId++
	} else {
		ordinalFrameId = 97
		pageNumber *= 10
	}

	if pageNumber > 999999999 {
		return "", errors.New("page number is greater than the maximum page number (999999999)")
	}

	return ConvertPidToPageId(pageNumber, string(ordinalFrameId))

}

func GetFollowOnPID(pageNumber int, frameId rune) (int, rune, error) {

	if !IsValidFrameId(frameId) || !IsValidPageNumber(pageNumber) {
		return 0, 0, errors.New("page numer or frame id is invalid")
	}

	// determine length of pageNumber
	if len(strconv.Itoa(pageNumber)) > 9 {
		return 0, 0, errors.New("page number is greater than the maximum page number (999999999)")
	}

	if frameId < 122 {
		frameId++
	} else {
		frameId = 97
		pageNumber *= 10
	}

	return pageNumber, frameId, nil

}

func ParseDataType(dataType string) (string, int) {

	var (
		arg int
		err error
	)

	// split the type as this allows for comma separated params.
	typ := strings.Split(dataType, ",")

	if len(typ) > 1 {
		dataType = typ[0]
		// check for number of rows
		if arg, err = strconv.Atoi(typ[1]); err == nil {
			if arg < 0 || arg > 22 {
				arg = 0
			}
		}
	}
	return dataType, arg
}

func IntToByte(i uint16) (byte, error) {

	var (
		byts []byte
		err  error
	)

	if byts, err = IntToBytes(i); err != nil {
		return 0, err
	}
	return byts[0], nil
}

func IntToBytes(i uint16) ([]byte, error) {

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func CheckPasswordStrength(password string) error {

	// entropy is a float64, representing the strength in base 2 (bits)
	//entropy := passwordvalidator.GetEntropy(password)
	//print(entropy)

	const minEntropyBits = 70 // seems to be a reasonable number
	// if the password has enough entropy, err is nil
	// otherwise, a formatted error message is provided explaining
	// how to increase the strength of the password
	// (safe to show to the client)
	return passwordvalidator.Validate(password, minEntropyBits)
}

func GetDateTimeFromUnixData(unixdate int64) string {

	var tm = time.Unix(unixdate, 0)

	// format using the magic date string - Mon Jan 2 15:04:05 MST 2006
	return tm.Format("2 January 2006 15:04")
}

func GetDateFromUnixData(unixdate int64) string {
	var tm = time.Unix(unixdate, 0)

	// format using the magic date string - Mon Jan 2 15:04:05 MST 2006
	return tm.Format("2 January 2006")
}

func GetTimeFromUnixData(unixdate int64) string {

	var tm = time.Unix(unixdate, 0)

	// format using the magic date string - Mon Jan 2 15:04:05 MST 2006
	return tm.Format("15:04")
}

func TruncateToStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func TruncateToEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func TruncateToStartOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

func TruncateToEndOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 0, t.Location())
}

func ToUpperCamelCase(word string) string {
	if len(word) > 2 {
		return strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	} else {
		return word
	}
}

func SetEvenParity(b byte) byte {

	v := b
	v = (v & 0x55) + ((v >> 1) & 0x55)
	v = (v & 0x33) + ((v >> 2) & 0x33)
	bits := (v + (v >> 4)) & 0xF

	if bits%2 > 0 {
		b = b | 0x80
	}
	return b
}

func FormatLogPreAmble(sessionCount int, connectionNumber int, ipAddress string) string {

	if len(ipAddress) == 0 {
		ipAddress = "0.0.0.0"
	}

	return fmt.Sprintf("%d:%d:%s: ", sessionCount, connectionNumber, ipAddress)
}

func GetIpAddress(conn net.Conn) string {

	var ipAddress string

	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		ipAddress = addr.IP.String()
	}

	return ipAddress
}
