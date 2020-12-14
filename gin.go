package spring

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/timeout"

	"github.com/gin-gonic/gin"
)

type key int

const (
	ginKey key = iota
)

type GinInitHandler func(ginEngine *gin.Engine)

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

func afterRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := getContainerFromRequest(c.Request.Context()).Delete()
		if err != nil {
			panic(err)
		}
	}
}

func InitGin(server graphql.ExecutableSchema, ginInitHandler GinInitHandler) *gin.Engine {
	if DIC().App().IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	} else if DIC().App().IsInTestMode() {
		gin.SetMode(gin.TestMode)
	}

	ginEngine := gin.New()

	if !DIC().App().IsInProdMode() {
		ginEngine.Use(gin.Logger())
	}

	ginEngine.Use(contextToContextMiddleware())
	ginEngine.Use(afterRequestMiddleware())

	ginEngine.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))
	ginEngine.GET("/", playgroundHandler())

	if ginInitHandler != nil {
		ginInitHandler(ginEngine)
	}

	return ginEngine
}

func graphqlHandler(server graphql.ExecutableSchema) gin.HandlerFunc {
	h := handler.New(server)

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New(1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		var message string
		asErr, is := err.(error)
		if is {
			message = asErr.Error()
		} else {
			message = "panic"
		}
		errorMessage := message + "\n" + string(debug.Stack())
		DIC().Log().Error(errorMessage)
		return errors.New("internal server error")
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
