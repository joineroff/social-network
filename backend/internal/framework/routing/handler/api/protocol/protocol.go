package protocol

import (
	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/usecase"
)

func Handle(ctx *gin.Context) (*Request, *Response, error) {
	req := NewRequest()
	res := NewResponse()

	if err := ctx.ShouldBindJSON(req); err != nil {
		res.WithError(NewError(usecase.NewValidationError("request structure mismath")))

		return req, res, err
	}

	if len(req.Meta) == 0 {
		req.Meta = make(map[string]interface{}, 0)
	}

	res.Meta = req.Meta

	return req, res, nil
}
