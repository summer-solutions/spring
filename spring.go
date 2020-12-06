package spring

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sarulabs/di"
)

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
	ginEngine := initGin(s, server)
	preDeploy()
	ginEngine.GET("/", playgroundHandler())

	panic(ginEngine.Run(":" + port))
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
	container = ioCBuilder.Build()
}
