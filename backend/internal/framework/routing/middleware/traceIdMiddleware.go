package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const TraceIDKey = "traceID"

func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID, err := uuid.NewRandom()
		if err != nil {
			zap.S().Errorw("failed to generate trace id", "error", err)
		}

		c.Set(TraceIDKey, traceID)
	}
}
