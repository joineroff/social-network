package protocol

type Response struct {
	Success int8                   `json:"success"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
	Error   *Error                 `json:"error,omitempty"`

	StatusCode int `json:"-"`
}

func NewResponse() *Response {
	return &Response{
		Success: 1,
	}
}

func (r *Response) WithError(e *Error) {
	r.Success = 0
	r.Error = e
	r.StatusCode = e.StatusCode
}
