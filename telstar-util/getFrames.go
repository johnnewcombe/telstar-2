package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

func cmdGetFrames(apiUrl string, saveDirectory string, primary bool, token string) (ResponseData, error) {

	var (
		respData ResponseData
		url      = apiUrl + "/frame"
		data     []byte
		result   []map[string]interface{}
		count    int
	)
	if primary {
		url += "?db=primary"
	}

	respData, err := get(url, token)
	if err != nil {
		return respData, err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}

	//parse the json array of tmp as unstructured data
	if err = json.Unmarshal([]byte(respData.Body), &result); err != nil {
		return respData, err
	}

	for _, frame := range result {

		// get the pid so that we can create the filename
		pid := frame["pid"].(map[string]interface{}) // this is a type assertion to convert from interface{}

		// get underlying data types
		pageNo := int(pid["page-no"].(float64))
		frameId := pid["frame-id"].(string)

		// convert each frame back to json so that each frame can be saved individually
		if data, err = json.MarshalIndent(frame, "", "    "); err != nil {
			return respData, err
		}

		// create filename
		filename := filepath.Join(saveDirectory, fmt.Sprintf("%d%s.json", pageNo, frameId))

		// save the file
		if err = saveText(filename, string(data)); err != nil {
			return respData, err
		}
		count++
	}

	respData.Body = fmt.Sprintf("{\"records saved\": %d}\r\n", count)
	return respData, nil
}
