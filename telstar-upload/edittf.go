package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/convert"
	"encoding/json"
	"errors"
	"strings"
)

func createEditTfFrame(pid Pid, urldata string) (frame string, err error) {

	var etfFrame EditTfFrame

	urldata = strings.TrimRight(urldata, "\n")
	urldata = strings.TrimRight(urldata, "\r")

	// check file for edit.tf url and create a simple edit.tf frame
	if isEditTfFrame(urldata) || isZxNetFrame(urldata) {
		// frame is edit.tf so create a json file
		//get the frame id from the filename
		etfFrame.PID = pid
		etfFrame.Content.Data = urldata
		etfFrame.Content.Type = "edit.tf"
		etfFrame.AuthorId = "telstar-util"
		etfFrame.StaticPage = true
		etfFrame.Visible = true
		etfFrame.FrameType = "information"
		etfFrame.RoutingTable = createDefaultRoutingTable(pid.PageNumber)

		j, err := json.Marshal(etfFrame)
		if err != nil {
			return "", err
		}
		return string(j), nil
	}
	return "", errors.New("validating data: edit.tf data not valid:")
}

func isEditTfFrame(urldata string) bool {

	if !(strings.HasPrefix(urldata, "https://edit.tf/#0:") || strings.HasPrefix(urldata, "http://edit.tf/#0:")) {
		return false
	}

	if _, err := convert.GetEditTfRawData(urldata); err != nil {
		return false
	}

	return true
}

func isZxNetFrame(urldata string) bool {

	if !(strings.HasPrefix(urldata, "https://zxnet.co.uk/teletext/editor/#0:") ||
		strings.HasPrefix(urldata, "http://zxnet.co.uk/teletext/editor/#0")) {
		return false
	}

	// we can use this original edit.tf function to get zxnet editor raw data
	if _, err := convert.GetEditTfRawData(urldata); err != nil {
		return false
	}

	return true
}
