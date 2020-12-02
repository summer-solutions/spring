package services

import (
	apexLog "github.com/apex/log"
	"github.com/sarulabs/di"
	diLocal "github.com/summer-solutions/spring/di"
	"github.com/summer-solutions/spring/services/log"
)

func LogGlobal(provider ...log.FieldProvider) *diLocal.ServiceDefinition {
	return &diLocal.ServiceDefinition{
		Name:   "log",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			l := apexLog.WithFields(&apexLog.Fields{})
			for _, fields := range provider {
				l = l.WithFields(fields())
			}
			return l, nil
		},
	}
}

func LogForRequest(provider ...log.RequestFieldProvider) *diLocal.ServiceDefinition {
	return &diLocal.ServiceDefinition{
		Name:   "log_request",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			l := log.New(provider...)
			return l, nil
		},
	}
}
