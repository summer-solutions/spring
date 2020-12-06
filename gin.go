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
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/gin-contrib/timeout"

	"github.com/gin-gonic/gin"
)

type key int

const (
	ginKey key = iota
)

func GinFromContext(ctx context.Context) *gin.Context {
	return ctx.Value(ginKey).(*gin.Context)
}

type GinMiddleWareProvider func() gin.HandlerFunc

func contextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func initGin(s *Spring, server graphql.ExecutableSchema) *gin.Engine {
	if App().IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	ginEngine := gin.New()

	if App().IsInProdMode() {
		h, has := GetServiceOptional("log_handler")
		if !has {
			log.SetHandler(h.(log.Handler))
		} else {
			log.SetHandler(json.Default)
		}
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
		ginEngine.Use(gin.Logger())
	}

	ginEngine.Use(contextToContextMiddleware())
	for _, provider := range s.middlewares {
		middleware := provider()
		if middleware != nil {
			ginEngine.Use(middleware)
		}
	}

	ginEngine.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))

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
		Log().Error(errorMessage)
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
