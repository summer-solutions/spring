package spring

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/sarulabs/di"
)

type Script interface {
	Code() string
	Description() string
	Run() error
	Unique() bool
}

type ScriptInterval interface {
	Script
	Interval() string
}

func (s *Spring) RunScript(script Script) {
	s.runScript(script)
}

func (s *Spring) RunScriptInterval(script ScriptInterval) {
	go func(script ScriptInterval) {
		for {
			valid := s.runScript(script)
			//TODO
			if valid {
				time.Sleep(time.Minute)
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}(script)
}

func ServiceDefinitionDynamicScript(scripts ...Script) *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "scripts",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return scripts, nil
		},
		Flags: func(registry *FlagsRegistry) {
			total := len(scripts)
			if len(scripts) > 0 {
				registry.Bool("scripts", false, fmt.Sprintf("list all %d available scripts", total))
				registry.String("run-script", "", "run script")
			}
		},
	}
}

func (s *Spring) runScript(script Script) bool {
	return func() bool {
		valid := true
		defer func() {
			if err := recover(); err != nil {
				var message string
				asErr, is := err.(error)
				if is {
					message = asErr.Error()
				} else {
					message = "panic"
				}
				DIC().Log().Error(message + "\n" + string(debug.Stack()))
				valid = false
			}
		}()
		err := script.Run()
		if err != nil {
			panic(err)
		}
		return valid
	}()
}
