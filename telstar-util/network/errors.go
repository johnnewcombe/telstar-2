package network

type RequestError struct {
}

func (e *RequestError) Error() string {
	return ""
}
