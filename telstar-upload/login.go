package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResponseData struct {
	Status     string
	StatusCode int
	Token      string
	Body       string
}

func cmdLogin(apiUrl string, userId string, password string) (ResponseData, error) {
	var (
		err      error
		respData ResponseData
	)

	//login and get the token
	if respData, err = login(apiUrl, userId, password); err != nil {
		return respData, err
	}

	return respData, nil
}

func login(apiUrl, userId, password string) (ResponseData, error) {

	var (
		respData ResponseData
	)

	// specific case of put that returns a token
	data := fmt.Sprintf("{\"user-id\": \"%s\", \"password\": \"%s\"}", userId, password)
	url := apiUrl + "/login"

	client := &http.Client{}

	bytData := []byte(data)

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytData))
	if err != nil {
		return respData, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return respData, err
	}

	cookies := resp.Cookies()
	if len(cookies) == 1 {
		respData.Token = cookies[0].Value
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}
	respData.StatusCode = resp.StatusCode
	respData.Status = resp.Status
	respData.Body = string(body)

	return respData, nil

}
