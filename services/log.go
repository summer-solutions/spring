package services

import (
	apexLog "github.com/apex/log"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/services/log"
)

func LogGlobal(provider ...log.FieldProvider) *spring.CDServiceDefinition {
	return &spring.CDServiceDefinition{
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

func LogForRequest(provider ...log.RequestFieldProvider) *spring.CDServiceDefinition {
	return &spring.CDServiceDefinition{
		Name:   "log_request",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			l := log.New(ctn, provider...)
			return l, nil
		},
	}
}
