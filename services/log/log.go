package log

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring/app"
)

type RequestFieldProvider func(ctx *gin.Context) log.Fielder
type FieldProvider func() log.Fielder

type RequestLog struct {
	providers []RequestFieldProvider
	entry     log.Interface
	container di.Container
}

func New(provider ...RequestFieldProvider) *RequestLog {
	return &RequestLog{providers: provider}
}

func (l *RequestLog) Log(ctx *gin.Context) log.Interface {
	if l.entry == nil {
		appName := l.container.Get("app").(*app.App).Name
		entry := log.WithFields(&log.Fields{"app": appName})
		for _, p := range l.providers {
			fields := p(ctx)
			if fields != nil {
				entry = entry.WithFields(fields)
			}
		}
		l.entry = entry
	}
	return l.entry
}
