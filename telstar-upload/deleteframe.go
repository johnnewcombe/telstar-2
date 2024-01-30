package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"fmt"
)

func cmdDeleteFrame(apiUrl, pageId string, primary bool, token string) (ResponseData, error) {

	var (
		respData ResponseData
		url      = apiUrl + "/frame/" + pageId
	)

	if primary {
		url += "?db=primary"
		url+= "&purge=true"
	} else {
		url+="?purge=true"
	}

	respData, err := delete(url, token)
	if err != nil {
		return respData, err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}
	logger.LogInfo.Printf("Deleting frame %s.", pageId )

	return respData, nil
}
