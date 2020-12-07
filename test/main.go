package main

import (
	"time"

	"github.com/summer-solutions/spring"
)

type testScript struct {
	name        string
	description string
	unique      bool
}

func (script *testScript) Run() error {
	return nil
}

func (script *testScript) Unique() bool {
	return script.unique
}

func (script *testScript) Code() string {
	return script.name
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
	r.RegisterDIService(spring.ServiceDefinitionDynamicScript(&testScript{"hello", "takie tam", false},
		&testScript{"hello-second", "takie tam inne description", true})).Build()
}
