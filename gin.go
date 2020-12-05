package spring

import (
	"context"

	"github.com/gin-gonic/gin"
)

type key int

const (
	ginKey key = iota
)

func GinFromContext(ctx context.Context) *gin.Context {
	return ctx.Value(ginKey).(*gin.Context)
}

func contextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
