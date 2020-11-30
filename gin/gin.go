package gin

import (
	"context"

	"github.com/gin-gonic/gin"
)

func FromContext(ctx context.Context) *gin.Context {
	return ctx.Value("GinContextKey").(*gin.Context)
}
