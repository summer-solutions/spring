package spring

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ryanuber/columnize"
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
				registry.Bool("list-scripts", false, fmt.Sprintf("list all %d available scripts", total))
				registry.String("run-script", "", "run script")
			}
		},
	}
}

func listScrips() {
	service, has := GetServiceOptional("scripts")
	if has {
		output := []string{
			"NAME | UNIQUE | DESCRIPTION ",
		}
		for _, def := range service.([]Script) {
			var unique string
			if def.Unique() {
				unique = "true"
			}
			output = append(output, strings.Join([]string{def.Code(), unique, def.Description()}, " | "))
		}
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
	}
	os.Exit(0)
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
