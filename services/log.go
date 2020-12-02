package services

import (
	apexLog "github.com/apex/log"
	"github.com/summer-solutions/spring/di"
	"github.com/summer-solutions/spring/services/log"
)

func LogGlobal(provider ...log.FieldProvider) *di.ServiceDefinition {
	return &di.ServiceDefinition{
		Name:   "log",
		Global: true,
		Build: func() (interface{}, error) {
			l := apexLog.WithFields(&apexLog.Fields{})
			for _, fields := range provider {
				l = l.WithFields(fields())
			}
			return l, nil
		},
	}
}

func LogForRequest(provider ...log.RequestFieldProvider) *di.ServiceDefinition {
	return &di.ServiceDefinition{
		Name:   "log_request",
		Global: false,
		Build: func() (interface{}, error) {
			l := log.New(provider...)
			return l, nil
		},
	}
}
