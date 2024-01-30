package api

import (
	"github.com/go-chi/render"
	"net/http"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "bad request",
		ErrorText:      err.Error(),
	}
}

func ErrTeapotRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 418,
		StatusText:     "I'm a teapot",
		ErrorText:      err.Error(),
	}
}

func ErrServerRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 500,
		StatusText:     "internal server error",
		ErrorText:      err.Error(),
	}
}

func ErrUnauthorizedRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 401,
		StatusText:     "unauthorised",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "error rendering response",
		ErrorText:      err.Error(),
	}
}

func ErrNotFoundRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "resource not found",
		ErrorText:      err.Error(),
	}
}

//var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "not found"}
//var ErrUnAuthorised = &ErrResponse{HTTPStatusCode: 401, StatusText: "unauthorized."}

type ResultResponse struct {
	HTTPStatusCode int    `json:"-"`
	ResultText     string `json:"result"` // user-level status message
}

func (s *ResultResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// this sets the Http Status Code in the response
	render.Status(r, s.HTTPStatusCode)
	return nil
}

func Result(httpResponseCode int, msg string) render.Renderer {
	return &ResultResponse{
		// the HttpStatusCode gets set in the HttpResponse
		// the result Text appears in the body
		HTTPStatusCode: httpResponseCode,
		ResultText:     msg,
	}
}

//func ResultUpdated(msg string) render.Renderer {
//	return &ResultResponse{
//		HTTPStatusCode: 202,
//		ResultText: msg,
//	}
//}

/* see https://httpstatuses.com/

1×× Informational

	100 Continue
	101 Switching Protocols
	102 Processing

2×× Success

	200 OK
	201 Created
	202 Accepted
	203 Non-authoritative Information
	204 No Content
	205 Reset Content
	206 Partial Content
	207 Multi-Status
	208 Already Reported
	226 IM Used

3×× Redirection

	300 Multiple Choices
	301 Moved Permanently
	302 Found
	303 See Other
	304 Not Modified
	305 Use Proxy
	307 Temporary Redirect
	308 Permanent Redirect

4×× Client Error

	400 Bad Request
	401 Unauthorized
	402 Payment Required
	403 Forbidden
	404 Not Found
	405 Method Not Allowed
	406 Not Acceptable
	407 Proxy Authentication Required
	408 Request Timeout
	409 Conflict
	410 Gone
	411 Length Required
	412 Precondition Failed
	413 Payload Too Large
	414 Request-URI Too Long
	415 Unsupported Media Type
	416 Requested Range Not Satisfiable
	417 Expectation Failed
	418 I'm a teapot
	421 Misdirected Request
	422 Unprocessable Entity
	423 Locked
	424 Failed Dependency
	426 Upgrade Required
	428 Precondition Required
	429 Too Many Requests
	431 Request Header Fields Too Large
	444 Connection Closed Without Response
	451 Unavailable For Legal Reasons
	499 Client Closed Request

5×× Server Error

	500 Internal Server Error
	501 Not Implemented
	502 Bad Gateway
	503 Service Unavailable
	504 Gateway Timeout
	505 HTTP Version Not Supported
	506 Variant Also Negotiates
	507 Insufficient Storage
	508 Loop Detected
	510 Not Extended
	511 Network Authentication Required
	599 Network Connect Timeout Error

*/
