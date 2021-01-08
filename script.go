package spring

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ryanuber/columnize"
)

type Script interface {
	Description() string
	Run(ctx context.Context, exit Exit)
	Unique() bool
}

type Exit interface {
	Valid()
	Error()
	Custom(exitCode int)
}

type exit struct {
	s *Spring
}

func (e *exit) Custom(exitCode int) {
	e.s.exit <- exitCode
}

func (e *exit) Valid() {
	e.Custom(0)
}

func (e *exit) Error() {
	e.Custom(1)
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
	go func() {
		for {
			valid := s.runScript(script)
			if !isInterval {
				s.done <- true
				break
			}
			//TODO
			if valid {
				time.Sleep(time.Minute)
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}()
	s.await()
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

func (s *Spring) runDynamicScrips(ctx context.Context, code string) {
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
			defScript.Run(ctx, &exit{s: s})
			return
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
		script.Run(s.ctx, &exit{s: s})
		return valid
	}()
}
