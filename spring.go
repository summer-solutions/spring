package spring

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/fatih/color"

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
	"github.com/sarulabs/di"
)

type GinMiddleWareProvider func() gin.HandlerFunc

type Spring struct {
	app                 *AppDefinition
	servicesDefinitions []*ServiceDefinition
	middlewares         []GinMiddleWareProvider
}

func New(appName string) *Spring {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeLocal
	}
	s := &Spring{app: &AppDefinition{Mode: mode, Name: appName}}
	return s
}

func (s *Spring) RegisterDIService(service ...*ServiceDefinition) *Spring {
	s.servicesDefinitions = append(s.servicesDefinitions, service...)
	return s
}

func (s *Spring) RegisterGinMiddleware(provider ...GinMiddleWareProvider) *Spring {
	s.middlewares = append(s.middlewares, provider...)
	return s
}

func (s *Spring) RunServer(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}

	s.initializeIoCHandlers()
	ginEngine := s.initGin(server)

	s.preDeploy()

	ginEngine.GET("/", playgroundHandler())

	panic(ginEngine.Run(":" + port))
}

func (s *Spring) preDeploy() {
	preDeployFlag := flag.Bool("pre-deploy", false, "Execute pre deploy mode")
	flag.Parse()

	if !*preDeployFlag {
		return
	}

	ormConfigService, has := OrmConfig()
	if !has {
		return
	}
	ormService := ormConfigService.CreateEngine()

	alters := ormService.GetAlters()

	hasAlters := false
	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
		hasAlters = true
	}

	if hasAlters {
		os.Exit(1)
	}

	if !s.app.IsInLocalMode() {
		os.Exit(0)
	}
}

func (s *Spring) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*ServiceDefinition{serviceLogGlobal(), serviceLogForRequest(), serviceConfig()}

	for _, def := range append(defaultDefinitions, s.servicesDefinitions...) {
		if def == nil {
			continue
		}

		var scope string
		if def.Global {
			scope = di.App
		} else {
			scope = di.Request
		}

		err := ioCBuilder.Add(di.Def{
			Name:  def.Name,
			Scope: scope,
			Build: def.Build,
			Close: def.Close,
		})
		if err != nil {
			panic(err)
		}
	}

	err := ioCBuilder.Add(di.Def{
		Name:  "app",
		Scope: di.App,
		Build: func(di.Container) (interface{}, error) {
			return s.app, nil
		},
	})

	if err != nil {
		panic(err)
	}
	setContainer(ioCBuilder.Build())
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

func (s *Spring) initGin(server graphql.ExecutableSchema) *gin.Engine {
	if s.app.IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	ginEngine := gin.New()

	if s.app.IsInProdMode() {
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

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
