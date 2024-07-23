package network

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

// ApiError create a simple struct for an error
type ApiError struct {
	msg string
}

func (e *ApiError) Error() string {
	//attach a message to the error
	return e.msg
}

func Get(apiUrl string, token string) (ResponseData, error) {

	var (
		respData ResponseData
	)
	client := &http.Client{}

	cookie := CreateCookie("token", token)
	jwtCookie := CreateCookie("jwt", token)

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

func Put(apiUrl string, data string, token string) (respData ResponseData, err error) {

	client := &http.Client{}
	cookie := CreateCookie("token", token)
	jwtCookie := CreateCookie("jwt", token)

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

	cookies := resp.Cookies()
	if len(cookies) == 1 {
		respData.Token = cookies[0].Value
	}

	respData.StatusCode = resp.StatusCode
	respData.Status = resp.Status
	respData.Body = string(body)

	return respData, nil
}

func Delete(apiUrl string, token string) (ResponseData, error) {

	client := &http.Client{}
	cookie := CreateCookie("token", token)
	jwtCookie := CreateCookie("jwt", token)

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

func CreateCookie(name string, token string) http.Cookie {

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
