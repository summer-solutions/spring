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
	"github.com/summer-solutions/orm"

	ginSpring "github.com/summer-solutions/spring/gin"

	"github.com/summer-solutions/spring/service"

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

type InitHandler func(s *Server, def *Def)
type GinMiddleware func(engine *gin.Engine) error
type Def struct {
	Name  string
	scope string
	Build func(ctn di.Container) (interface{}, error)
	Close func(obj interface{}) error
}

type Server struct {
	mode            string
	initHandlers    []InitHandler
	requestServices []InitHandler
	middlewares     []GinMiddleware
}

func NewServer(handler InitHandler, middlewares ...GinMiddleware) *Server {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeProd
	}

	s := &Server{mode: mode, middlewares: middlewares}

	s.initializeIoCHandlers(handler)

	s.preDeploy()
	return s
}

func (s *Server) Run(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	if s.IsInProdMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	if s.IsInProdMode() {
		h, err := service.GetGlobalContainer().SafeGet("log_handler")
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

	s.attachMiddlewares(r)

	r.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))
	r.GET("/", playgroundHandler())
	panic(r.Run(":" + port))
}

func (s *Server) RegisterGlobalServices(handlers ...InitHandler) {
	s.initHandlers = append(s.initHandlers, handlers...)
}

func (s *Server) RegisterRequestServices(handlers ...InitHandler) {
	s.requestServices = append(s.requestServices, handlers...)
}

func (s *Server) preDeploy() {
	preDeployFlag := flag.Bool("pre-deploy", false, "Execute pre deploy mode")
	flag.Parse()

	if !*preDeployFlag && !s.IsInLocalMode() {
		return
	}

	ormConfigService := service.OrmConfig().(orm.ValidatedRegistry)
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

func (s *Server) initializeIoCHandlers(handlerRegister InitHandler) {
	ioCBuilder, _ := di.NewBuilder()

	handlerRegister(s, nil)

	scopes := map[string][]InitHandler{di.App: s.initHandlers, di.Request: s.requestServices}
	for scope, services := range scopes {
		for _, callback := range services {
			def := &Def{scope: scope}

			callback(s, def)
			if def.Name == "" {
				panic("IoC " + scope + " service is registered without name")
			}

			if def.Build == nil {
				panic("IoC " + scope + " service is registered without Build function")
			}

			err := ioCBuilder.Add(di.Def{
				Name:  def.Name,
				Scope: def.scope,
				Build: def.Build,
				Close: def.Close,
			})
			if err != nil {
				panic(err)
			}
		}
	}
	service.SetGlobalContainer(ioCBuilder.Build())
}

func (s *Server) attachMiddlewares(engine *gin.Engine) {
	for _, middleware := range s.middlewares {
		err := middleware(engine)
		if err != nil {
			panic(err)
		}
	}
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
		service.Log().Error(message + "\n" + string(debug.Stack()))
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
