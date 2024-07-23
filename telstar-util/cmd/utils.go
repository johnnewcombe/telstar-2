package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func saveText(filename string, text string) error {
	d1 := []byte(text)
	err := ioutil.WriteFile(filename, d1, 0644)
	return err
}

func loadText(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	return string(dat), err
}

func parseFrame(jsonFrame string) (types.Frame, error) {

	var frame types.Frame
	frameBytes := []byte(jsonFrame)

	if !json.Valid(frameBytes) {
		return frame, fmt.Errorf("validating frame: invalid json")
	}

	if err := json.Unmarshal(frameBytes, &frame); err != nil {
		return frame, fmt.Errorf("parsing json: invalid")
	}

	return frame, nil
}

func CreateDefaultRoutingTable(pageNo int) []int {

	var (
		pageNumber float64
	)
	routingTable := make([]int, 11)

	// sort out the entries 0-9 (i.e.keys presses 0-9)
	for n := 0; n < 10; n++ {
		routingTable[n] = n + (pageNo * 10)
	}

	//sort out the hash route
	pageNumber = float64(pageNo)

	for pageNumber > 999 {
		pageNumber = math.Floor(pageNumber / 10)
	}
	routingTable[10] = int(pageNumber)

	return routingTable
}

func GetPidFromFileName(filename string) (pid types.Pid, ok bool) {

	var (
		regex = regexp.MustCompile("^[0-9]{1,10}[a-z].edit.tf$")
		pgNo  int
		err   error
	)
	_, filename = filepath.Split(filename)

	if regex.MatchString(filename) {

		pns := filename[0 : len(filename)-9]
		pnsl := len(pns)
		if pgNo, err = strconv.Atoi(pns); err != nil {
			return pid, false
		}
		pid.PageNumber = pgNo
		pid.FrameId = filename[pnsl : pnsl+1]

		return pid, true
	}

	return pid, false
}

func stdOut(cmd *cobra.Command, respData network.ResponseData, result map[string]string) {

	const (
		ResponseJson = `{"HTTP Status" : "%s", "Result" : %s}`
		ResponseText = "HTTP Status: %s" // result (if it exists, and NL added later
	)

	if cmd.Flags().Lookup("json").Changed {
		jsonData, _ := json.Marshal(result)
		fmt.Printf(ResponseJson+"\n\n", respData.Status, string(jsonData))

	} else {
		text := strings.Builder{}
		text.WriteString(fmt.Sprintf(ResponseText, respData.Status))

		if result != nil {
			text.WriteString(", ")
		}

		for key, value := range result {
			if len(value) == 0 {
				value = "null"
			}
			text.WriteString(key + ": " + value + ", ")
		}
		textS := text.String()

		// format the result if there is one
		if result != nil {
			textS = textS[:len(textS)-2] // remove comma and trailing space
		}

		fmt.Printf("%s\n\n", textS)
	}
}
