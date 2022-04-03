package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, struct{}{})
	}
}
