package spring

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ryanuber/columnize"
)

type Script interface {
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

type ScriptIntermediate interface {
	IsIntermediate() bool
}

type ScriptOptional interface {
	Active() bool
}

func (s *Spring) RunScript(script Script) {
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
}

func listScrips() {
	scripts := DIC().App().registry.scripts
	if len(scripts) > 0 {
		output := []string{
			"NAME | OPTIONS | DESCRIPTION ",
		}
		for _, defCode := range scripts {
			def := GetServiceRequired(defCode).(Script)
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
			intermediate, is := def.(ScriptIntermediate)
			if is && intermediate.IsIntermediate() {
				options = append(options, "intermediate")
			}
			output = append(output, strings.Join([]string{defCode, strings.Join(options, ","), def.Description()}, " | "))
		}
		_, _ = os.Stdout.WriteString(columnize.SimpleFormat(output) + "\n")
	}
}

func runDynamicScrips(code string) {
	scripts := DIC().App().registry.scripts
	if len(scripts) == 0 {
		panic(fmt.Sprintf("unknown script %s", code))
	}
	for _, defCode := range scripts {
		if defCode == code {
			def, has := GetServiceOptional(defCode)
			if !has {
				panic(fmt.Sprintf("unknown script %s", code))
			}
			defScript := def.(Script)
			err := defScript.Run()
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
