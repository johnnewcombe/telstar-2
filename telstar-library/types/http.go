package types

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type ApiResponse struct {
	HTTPStatusCode int    `json:"-"`
	ResultText     string `json:"result"` // user-level status message
}

func (s *ApiResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// this sets the Http Status Code in the response
	render.Status(r, s.HTTPStatusCode)
	return nil
}

func (s *ApiResponse) Load(jsonResponse string) error {

	var r ApiResponse
	jsonBytes := []byte(jsonResponse)

	if !json.Valid(jsonBytes) {

		return fmt.Errorf("validating frame: invalid json")
	}

	if err := json.Unmarshal(jsonBytes, &r); err != nil {
		return fmt.Errorf("parsing json: invalid")
	}
	s.ResultText = r.ResultText
	s.HTTPStatusCode = r.HTTPStatusCode

	return nil
}
