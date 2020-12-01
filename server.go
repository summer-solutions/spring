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

const ModeLocal = "local"
const ModeDev = "dev"
const ModeProd = "prod"

type Server struct {
	mode                  string
	cdServicesDefinitions []*CDServiceDefinition
	middlewares           []gin.HandlerFunc
}

func NewServer() *Server {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeLocal
	}
	s := &Server{mode: mode}
	return s
}

func (s *Server) RegisterCDService(service ...*CDServiceDefinition) *Server {
	s.cdServicesDefinitions = append(s.cdServicesDefinitions, service...)
	return s
}

func (s *Server) RegisterGinMiddleware(middleware ...gin.HandlerFunc) *Server {
	s.middlewares = append(s.middlewares, middleware...)
	return s
}

func (s *Server) Run(defaultPort uint, server graphql.ExecutableSchema) {
	s.preDeploy()
	s.initializeIoCHandlers()
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	if s.IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	if s.IsInProdMode() {
		h, err := GetContainer().SafeGet("log_handler")
		if err == nil {
			log.SetHandler(h.(log.Handler))
		} else {
			log.SetHandler(json.Default)
		}
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
		r.Use(gin.Logger())
	}

	r.Use(ginSpring.ContextToContextMiddleware())
	for _, m := range s.middlewares {
		if m != nil {
			r.Use(m)
		}
	}

	r.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))
	r.GET("/", playgroundHandler())
	panic(r.Run(":" + port))
}

func (s *Server) preDeploy() {
	preDeployFlag := flag.Bool("pre-deploy", false, "Execute pre deploy mode")
	flag.Parse()

	if !*preDeployFlag && !s.IsInLocalMode() {
		return
	}

	ormConfigService, has := CDOrmConfig()
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

	if !s.IsInLocalMode() {
		os.Exit(0)
	}
}

func (s *Server) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	for _, def := range s.cdServicesDefinitions {
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
	container = ioCBuilder.Build()
}

func (s *Server) IsInLocalMode() bool {
	return s.IsInMode(ModeLocal)
}

func (s *Server) IsInProdMode() bool {
	return s.IsInMode(ModeProd)
}

func (s *Server) IsInDevMode() bool {
	return s.IsInMode(ModeDev)
}

func (s *Server) IsInMode(mode string) bool {
	return s.mode == mode
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
		l, has := CDLog()
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
