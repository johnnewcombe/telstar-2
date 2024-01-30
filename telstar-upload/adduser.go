package main

import (
	"fmt"
)

func cmdAddUser(apiUrl, userId, password, token string) (ResponseData, error) {

	var (
		respData ResponseData
		url      = apiUrl + "/user"
	)

	// TODO the error should come from the api
	//if !isValidUserId(userId) {
	//	return errors.New("bad userid")
	//}

	data := "{\"user-id\":" + userId + ",\"password\": \"" + password + "\"}"
	respData, err := put(url, data, token)
	if err != nil {
		return respData, err
	}

	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return respData, fmt.Errorf("%s", respData.Body)
	}

	return respData, nil

}
