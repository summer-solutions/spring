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

func (r *Registry) Build() *Spring {
	r.initializeIoCHandlers()
	r.initializeLog()
	flags := DIC().App().Flags()
	if flags.Bool("list-scripts") {
		listScrips()
	}
	scriptToRun := flags.String("run-script")
	if scriptToRun != "" {
		runDynamicScrips(scriptToRun)
	}
	return &Spring{registry: r}
}

func (r *Registry) RegisterDIService(service ...*ServiceDefinition) *Registry {
	r.servicesDefinitions = append(r.servicesDefinitions, service...)
	return r
}

func (r *Registry) RegisterGinMiddleware(provider ...GinMiddleWareProvider) *Registry {
	r.middlewares = append(r.middlewares, provider...)
	return r
}

func (r *Registry) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*ServiceDefinition{
		serviceApp(r.app),
		serviceLogGlobal(),
		serviceLogForRequest(),
		serviceConfig(),
	}

	flagsRegistry := &FlagsRegistry{flags: make(map[string]interface{})}
	for _, def := range append(defaultDefinitions, r.servicesDefinitions...) {
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
	r.app.flags = &Flags{flagsRegistry}
}
