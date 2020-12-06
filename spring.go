package spring

import (
	"flag"
	"os"

	"github.com/sarulabs/di"
)

type Registry struct {
	app                 *AppDefinition
	servicesDefinitions []*ServiceDefinition
	middlewares         []GinMiddleWareProvider
}

type Spring struct {
	registry *Registry
}

func New(appName string) *Registry {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeLocal
	}
	return &Registry{app: &AppDefinition{mode: mode, name: appName}}
}

func (s *Registry) Build() *Spring {
	s.initializeIoCHandlers()
	s.initializeLog()
	return &Spring{registry: s}
}

func (s *Registry) RegisterDIService(service ...*ServiceDefinition) *Registry {
	s.servicesDefinitions = append(s.servicesDefinitions, service...)
	return s
}

func (s *Registry) RegisterGinMiddleware(provider ...GinMiddleWareProvider) *Registry {
	s.middlewares = append(s.middlewares, provider...)
	return s
}

func (s *Registry) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*ServiceDefinition{
		serviceApp(s.app),
		serviceLogGlobal(),
		serviceLogForRequest(),
		serviceConfig(),
	}

	flagsRegistry := &FlagsRegistry{flags: make(map[string]interface{})}
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
		if def.Flags != nil {
			def.Flags(flagsRegistry)
		}
	}

	err := ioCBuilder.Add()

	if err != nil {
		panic(err)
	}
	container = ioCBuilder.Build()
	dicInstance = &dic{}
	flag.Parse()
	s.app.flags = &Flags{flagsRegistry}
}
