package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, struct{}{})
	}
}
