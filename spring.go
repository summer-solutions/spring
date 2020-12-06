package spring

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

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
	s := &Spring{app: &AppDefinition{mode: mode, name: appName}}
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
	dicInstance = &dic{}
}

func (s *Spring) initializeLog() {
	if DIC().App().IsInProdMode() {
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
	}
}
