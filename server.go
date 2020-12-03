package spring

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/summer-solutions/spring/app"

	"github.com/summer-solutions/spring/ioc"

	"github.com/fatih/color"
	ginSpring "github.com/summer-solutions/spring/gin"

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

type Server struct {
	app                 *app.App
	servicesDefinitions []*ioc.ServiceDefinition
	middlewares         []GinMiddleWareProvider
}

func NewServer(appName string) *Server {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = app.ModeLocal
	}
	s := &Server{app: &app.App{Mode: mode, Name: appName}}
	return s
}

func (s *Server) RegisterDIService(service ...*ioc.ServiceDefinition) *Server {
	s.servicesDefinitions = append(s.servicesDefinitions, service...)
	return s
}

func (s *Server) RegisterGinMiddleware(provider ...GinMiddleWareProvider) *Server {
	s.middlewares = append(s.middlewares, provider...)
	return s
}

func (s *Server) InitGin(server graphql.ExecutableSchema) *gin.Engine {
	s.initializeIoCHandlers()

	if s.app.IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	ginEngine := gin.New()

	if s.app.IsInProdMode() {
		h, has := ioc.GetServiceOptional("log_handler")
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

	ginEngine.Use(ginSpring.ContextToContextMiddleware())
	for _, provider := range s.middlewares {
		middleware := provider()
		if middleware != nil {
			ginEngine.Use(middleware)
		}
	}

	ginEngine.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))

	return ginEngine
}

func (s *Server) Run(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}

	ginEngine := s.InitGin(server)

	s.preDeploy()

	ginEngine.GET("/", playgroundHandler())

	panic(ginEngine.Run(":" + port))
}

func (s *Server) preDeploy() {
	preDeployFlag := flag.Bool("pre-deploy", false, "Execute pre deploy mode")
	flag.Parse()

	if !*preDeployFlag && !s.app.IsInLocalMode() {
		return
	}

	ormConfigService, has := ioc.OrmConfig()
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

func (s *Server) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	for _, def := range s.servicesDefinitions {
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
	ioc.SetContainer(ioCBuilder.Build())
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
		l, has := ioc.Log()
		if has {
			l.Error(errorMessage)
		} else {
			log.Error(errorMessage)
		}
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
