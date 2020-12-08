package main

import (
	"time"

	"github.com/sarulabs/di"

	"github.com/summer-solutions/spring"
)

type testScript struct {
	description string
	unique      bool
}

func (script *testScript) Run() error {
	return nil
}

func (script *testScript) Unique() bool {
	return script.unique
}

func (script *testScript) Description() string {
	return script.description
}

func (script *testScript) Active() bool {
	return true
}

func (script *testScript) Interval() time.Duration {
	return 3 * time.Second
}

func main() {
	r := spring.New("test_script")
	r.RegisterDIService(&spring.ServiceDefinition{
		Name:   "aa",
		Global: true,
		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &testScript{"takie tam", false}, nil
		},
	})
	r.RegisterDIService(&spring.ServiceDefinition{
		Name:   "bb",
		Global: true,
		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &testScript{"takie tam dwa", true}, nil
		},
	})
	r.RegisterDIService(spring.ServiceProviderConfigDirectory("../config"))
	r.RegisterDIService().Build()
}
