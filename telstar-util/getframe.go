package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"errors"
	"fmt"
)

func cmdGetFrame(apiUrl string, pageId string, primary bool, token string) (ResponseData, error) {

	var (
		respData ResponseData
		err      error
	)

	if !utils.IsValidPageId(pageId) {
		return respData, errors.New("invalid frame id")
	}

	var url = apiUrl + "/frame/" + pageId
	if primary {
		url += "?db=primary"
	}

	respData, err = get(url, token)
	if err != nil {
		return respData, err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}

	return respData, nil

}
