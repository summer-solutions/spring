package log

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

type RequestFieldProvider func(ctx *gin.Context) log.Fielder
type FieldProvider func() log.Fielder

type RequestLog struct {
	providers []RequestFieldProvider
	entry     log.Interface
}

func New(provider ...RequestFieldProvider) *RequestLog {
	return &RequestLog{providers: provider}
}

func (l *RequestLog) Log(ctx *gin.Context) log.Interface {
	if l.entry == nil {
		entry := log.WithFields(&log.Fields{})
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
