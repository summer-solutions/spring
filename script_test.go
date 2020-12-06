package spring

import (
	"testing"

	"github.com/sarulabs/di"
	"github.com/tj/assert"
)

type testScript struct {
	RunCounter int
}

func (script *testScript) Run() error {
	script.RunCounter++
	return nil
}

func (script *testScript) Unique() bool {
	return false
}

func TestRunScript(t *testing.T) {
	s := New("test_script").RegisterDIService()
	testService := &ServiceDefinition{
		Name:   "test_service",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return "hello", nil
		},
	}
	s.RegisterDIService(testService)

	testScript := &testScript{}
	s.RunScript(testScript)
	assert.Equal(t, 1, testScript.RunCounter)
	assert.Equal(t, "test_script", DIC().App().Name())
	assert.Equal(t, "hello", GetServiceRequired("test_service"))
}
