package network

type ResponseData struct {
	Status     string
	StatusCode int
	Token      string
	Body       string
}

func (r *ResponseData) SetOK() {
	r.Status = "OK"
	r.StatusCode = 200
}
