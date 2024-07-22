package main

/*

import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"fmt"
)

func cmdPublishFrame(apiUrl string, pageId string, primary bool, token string) (ResponseData, error) {

	var (
		respData ResponseData
		err      error
		url      = apiUrl + "/publish/" + pageId
	)

	if primary {
		url += "?db=primary"
	}

	if !utils.IsValidPageId(pageId) {
		exitWithHelp()
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

*/
