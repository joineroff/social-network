package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/dto"
	"github.com/joineroff/social-network/backend/internal/framework/routing/handler/api/protocol"
	"github.com/joineroff/social-network/backend/internal/usecase"
)

func SignInHandler(uc *usecase.SignInUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, res, err := protocol.Handle(c)
		if err != nil {
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		signInDto := &dto.SignInInputDto{}

		if req.Data == nil {
			e := usecase.NewValidationError("empty data provided")
			res.WithError(protocol.NewError(e))
			c.JSON(res.StatusCode, res)

			return
		}

		err = json.Unmarshal(*req.Data, signInDto)
		if err != nil {
			e := usecase.NewValidationError(err.Error())
			res.WithError(protocol.NewError(e))
			c.JSON(res.StatusCode, res)

			return
		}

		output, err := uc.Do(c, signInDto)
		if err != nil {
			res.WithError(protocol.NewError(err))
			c.JSON(res.StatusCode, res)

			return
		}

		c.JSON(http.StatusOK, output)
	}
}
