package main

import (
	"fmt"
)

func cmdDeleteUser(apiUrl, userId string, token string) (ResponseData, error) {

	var (
		respData ResponseData
		url      = apiUrl + "/user/" + userId
	)

	respData, err := delete(url, token)
	if err != nil {
		return respData, err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}

	return respData, nil

}
