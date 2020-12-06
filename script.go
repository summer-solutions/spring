package spring

import (
	"runtime/debug"
	"time"
)

type Script interface {
	Run() error
	Unique() bool
}

type ScriptInterval interface {
	Script
	Interval() string
}

func (s *Spring) RunScript(script Script) {
	s.initializeIoCHandlers()
	s.initializeLog()
	s.runScript(script)
}

func (s *Spring) RunScriptInterval(script ScriptInterval) {
	s.initializeIoCHandlers()
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

func (s *Spring) runScript(script Script) bool {
	s.initializeIoCHandlers()
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
