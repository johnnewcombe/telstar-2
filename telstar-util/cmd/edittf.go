package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"encoding/json"
	"errors"
	"strings"
)

type EditTfFrame struct {
	PID          types.Pid     `json:"pid" bson:"pid"`
	Visible      bool          `json:"visible" bson:"visible"`
	Content      types.Content `json:"content" bson:"content"`
	RoutingTable []int         `json:"routing-table" bson:"routing-table"`
	AuthorId     string        `json:"author-id" bson:"author-id"`
	StaticPage   bool          `json:"static-page" bson:"static-page"`
}

// FIXME No longer supported
func createEditTfFrame(pid types.Pid, urldata string) (frame string, err error) {

	var etfFrame EditTfFrame

	urldata = strings.TrimRight(urldata, "\n")
	urldata = strings.TrimRight(urldata, "\r")

	// check file for edit.tf url and create a simple edit.tf frame
	if isEditTfFrame(urldata) {
		// frame is edit.tf so create a json file
		//get the frame id from the filename
		etfFrame.PID = pid
		etfFrame.Content.Data = urldata
		etfFrame.Content.Type = "edit.tf"
		etfFrame.AuthorId = "telstar-util"
		etfFrame.StaticPage = true
		etfFrame.Visible = true
		etfFrame.RoutingTable = CreateDefaultRoutingTable(pid.PageNumber)

		j, err := json.Marshal(etfFrame)
		if err != nil {
			return "", err
		}
		return string(j), nil
	}
	return "", errors.New("validating data: edit.tf data not valid:")
}

func isEditTfFrame(urldata string) bool {

	if !(strings.HasPrefix(urldata, "https://edit.tf/#0:") ||
		strings.HasPrefix(urldata, "http://edit.tf/#0:")) {
		return false
	}
	// TODO Check the length also
	// hash position +3 is start of data len is 1120 or 1167
	data := strings.Split(urldata, "#0:")
	if len(data) != 2 {
		return false
	}
	rawEdata := strings.TrimRight(data[1], "\n")
	rawEdata = strings.TrimRight(rawEdata, "\r")

	length := len(rawEdata)
	if !(length == 1120 || length == 1167) {
		return false
	}

	return true
}
