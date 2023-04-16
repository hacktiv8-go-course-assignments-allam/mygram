package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mygram/go-common/pkg/context"
)

func CorrelationIDInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// try get from header
		val := ctx.GetHeader(context.CorrID.String())
		if val == "" {
			val = uuid.New().String()
		}

		ctx.Set(context.CorrID.String(), val)
		ctx.Next()
	}
}
