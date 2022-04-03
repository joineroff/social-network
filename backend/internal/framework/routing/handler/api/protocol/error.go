package protocol

import (
	"net/http"

	"github.com/joineroff/social-network/backend/internal/usecase"
)

type Error struct {
	StatusCode int               `json:"-"`
	Type       string            `json:"type"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details"`
}

func NewError(err error) *Error {
	e := &Error{}

	switch cast := err.(type) {
	case usecase.ValidationError:
		e.StatusCode = http.StatusBadRequest
		e.Type = "validation error"
		e.Message = cast.Message

		if len(cast.Details) > 0 {
			e.Details = cast.Details
		}
	case usecase.LogicError:
		e.StatusCode = http.StatusUnprocessableEntity
		e.Type = "logic error"
		e.Message = cast.Message

		if len(cast.Details) > 0 {
			e.Details = cast.Details
		}
	default:
		e.StatusCode = http.StatusInternalServerError
		e.Type = "internal error"
		e.Message = "something went wrong"
	}

	return e
}
