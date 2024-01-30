package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type ApiJsonError struct {
	Data string `json:"error" bson:"error"`
}

type HTTPErrorResponse struct {
	Error string `json:"error" bson:"error"`
}

// create a simple struct for an error
type ApiError struct {
	msg string
}

func (e *ApiError) Error() string {
	//attach a message to the error
	return e.msg
}

func get(apiUrl string, token string) (ResponseData, error) {

	var (
		respData ResponseData
	)
	client := &http.Client{}

	cookie := createCookie("token", token)
	jwtCookie := createCookie("jwt", token)

	var bytData []byte

	req, err := http.NewRequest(http.MethodGet, apiUrl, bytes.NewBuffer(bytData))
	if err != nil {
		return respData, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if len(token) > 0 {
		req.AddCookie(&cookie)
		req.AddCookie(&jwtCookie)
	}
	resp, err := client.Do(req)

	if err != nil {
		return respData, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}
	respData.StatusCode = resp.StatusCode
	respData.Status = resp.Status
	respData.Body = string(body)

	return respData, nil
}

func put(apiUrl string, data string, token string) (respData ResponseData, err error) {

	client := &http.Client{}
	cookie := createCookie("token", token)
	jwtCookie := createCookie("jwt", token)

	bytData := []byte(data)

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(bytData))
	if err != nil {
		return respData, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if len(token) > 0 {
		req.AddCookie(&cookie)
		req.AddCookie(&jwtCookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return respData, err
	}

	defer req.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}

	respData.StatusCode = resp.StatusCode
	respData.Status = resp.Status
	respData.Body = string(body)

	return respData, nil
}

func delete(apiUrl string, token string) (ResponseData, error) {

	client := &http.Client{}
	cookie := createCookie("token", token)
	jwtCookie := createCookie("jwt", token)

	var (
		bytData  []byte
		respData ResponseData
	)

	req, err := http.NewRequest(http.MethodDelete, apiUrl, bytes.NewBuffer(bytData))
	if err != nil {
		return respData, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if len(token) > 0 {
		req.AddCookie(&cookie)
		req.AddCookie(&jwtCookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return respData, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}

	respData.StatusCode = resp.StatusCode
	respData.Status = resp.Status
	respData.Body = string(body)

	return respData, nil
}

func purge(apiUrl string, token string) (ResponseData, error) {

	client := &http.Client{}
	cookie := createCookie("token", token)
	jwtCookie := createCookie("jwt", token)

	var (
		bytData  []byte
		respData ResponseData
	)

	req, err := http.NewRequest(http.MethodDelete, apiUrl, bytes.NewBuffer(bytData))
	if err != nil {
		return respData, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	if len(token) > 0 {
		req.AddCookie(&cookie)
		req.AddCookie(&jwtCookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return respData, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}

	respData.StatusCode = resp.StatusCode
	respData.Status = resp.Status
	respData.Body = string(body)

	return respData, nil
}
/*
func post(apiUrl string) (*http.Response, error) {

	var buf io.Reader
	var resp *http.Response

	defer resp.Body.Close()

	resp, err := http.Post(apiUrl, "image/jpeg", &buf)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
*/

func createCookie(name string, token string) http.Cookie {

	expire := time.Now().AddDate(0, 0, 1)
	return http.Cookie{
		Name:       name,
		Value:      token,
		Path:       "",
		Domain:     "",
		Expires:    expire,
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   []string{},
	}
}
func saveText(filename string, text string) error {
	d1 := []byte(text)
	err := ioutil.WriteFile(filename, d1, 0644)
	return err
}
func loadText(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	return string(dat), err
}

/*
func parseHttpError(httpResponse string) error {
//FIXME does this need to respond to HTTP STATUS CODES?? or should api 2.0 follow the lines here
	var (
		httpErrorResponse HTTPErrorResponse
	)

	if err := json.Unmarshal([]byte(httpResponse), &httpErrorResponse); err != nil {
			return nil // unmarshal errors mean that it cant be an http error
	}
	errValue := strings.ToLower(httpErrorResponse.Error)

	if len(errValue) > 0 {
		return fmt.Errorf("{http error: %s}", errValue)
	}

	// TODO only for version 2.0 at the moment
	//if statusValue != "OK"{
	//	return fmt.Errorf("http error: %s", errValue)
	//}
	return nil
}
*/
/*
//TODO remove the need for this by using parseHttpError see getframe
func parseError(jsonResponse string) error {

	var apiError ApiJsonError
	responseBytes := []byte(jsonResponse)
	if err := json.Unmarshal(responseBytes, &apiError); err != nil {
		// if we cannot marshal the response then it is not an error
		return nil
	}
	// we get here if we could marshal the response against the ApiError type
	// but it could be an "ok" nessage so check it
	switch strings.ToLower(apiError.Data) {
	case "forbidden":
		return &ApiError{
			msg: "{\"error\":\"log on required\"}",
		}
	case "not acceptable":
		return &ApiError{
			msg: "{\"error\":\"Frame not acceptable\"}",
		}
	case "login":

	}
	// this is good i.e. no json error message
	return nil

}
*/
