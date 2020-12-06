package main

import (
	"github.com/summer-solutions/spring"
)

type testScript struct {
}

func (script *testScript) Run() error {
	return nil
}

func (script *testScript) Unique() bool {
	return false
}

func (script *testScript) Code() string {
	return "test script"
}

func (script *testScript) Description() string {
	return "test description"
}

func main() {
	r := spring.New("test_script").RegisterDIService()
	r.RegisterDIService(spring.ServiceDefinitionDynamicScript(&testScript{})).Build()
}
