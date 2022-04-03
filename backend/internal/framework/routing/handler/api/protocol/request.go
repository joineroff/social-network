package protocol

import "encoding/json"

type Request struct {
	Meta map[string]interface{} `json:"meta,omitempty"`
	Data *json.RawMessage       `json:"data,omitempty"`
}

func NewRequest() *Request {
	return &Request{}
}
