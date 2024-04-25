package api

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	Success = "OK"
	Error   = "Error"
)

func SuccessResponse() *Response {
	return &Response{Status: Success}
}
func ErrorResponse(msg string) *Response {
	return &Response{Status: Error, Error: msg}
}
