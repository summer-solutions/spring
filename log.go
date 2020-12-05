package spring

import (
	apexLog "github.com/apex/log"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring/ioc"
	"github.com/summer-solutions/spring/services/log"
)

func logGlobal() *ioc.ServiceDefinition {
	return &ioc.ServiceDefinition{
		Name:   "log",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			l := apexLog.WithFields(&apexLog.Fields{"app": ioc.App().Name})
			key := "_log_providers"
			_, has := ctn.Definitions()[key]
			if has {
				providers := ctn.Get(key).([]log.FieldProvider)
				for _, fields := range providers {
					l = l.WithFields(fields())
				}
			}
			return l, nil
		},
	}
}

func logForRequest() *ioc.ServiceDefinition {
	return &ioc.ServiceDefinition{
		Name:   "log_request",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			key := "_log_request_providers"
			_, has := ctn.Definitions()[key]
			var l *log.RequestLog
			if has {
				providers := ctn.Get(key).([]log.RequestFieldProvider)
				l = log.New(providers...)
			} else {
				l = log.New()
			}
			return l, nil
		},
	}
}
