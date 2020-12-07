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
	Interval() time.Duration
}

type ScriptIntervalOptional interface {
	IntervalActive() bool
}

type ScriptOptional interface {
	Active() bool
}

func (s *Spring) RunScript(script Script) {
	go func(script Script) {
		_, isInterval := script.(ScriptInterval)
		for {
			valid := s.runScript(script)
			if !isInterval {
				break
			}
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
			"NAME | OPTIONS | DESCRIPTION ",
		}
		for _, def := range service.([]Script) {
			options := make([]string, 0)
			interval, is := def.(ScriptInterval)
			if is {
				options = append(options, "interval")
				duration := "every " + interval.Interval().String()
				_, is := def.(ScriptIntervalOptional)
				if is {
					duration += " with condition"
				}
				options = append(options, duration)
			}

			if def.Unique() {
				options = append(options, "unique")
			}
			optional, is := def.(ScriptOptional)
			if is {
				options = append(options, "optional")
				if optional.Active() {
					options = append(options, "active")
				} else {
					options = append(options, "inactive")
				}
			}
			output = append(output, strings.Join([]string{def.Code(), strings.Join(options, ","), def.Description()}, " | "))
		}
		_, _ = os.Stdout.WriteString(columnize.SimpleFormat(output))
	}
}

func runDynamicScrips(code string) {
	service, has := GetServiceOptional("scripts")
	if !has {
		panic(fmt.Sprintf("unknown script %s", code))
	}
	for _, def := range service.([]Script) {
		if def.Code() == code {
			err := def.Run()
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	}
	panic(fmt.Sprintf("unknown script %s", code))
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
