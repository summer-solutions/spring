package log

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
)

type RequestFieldProvider func(ctx *gin.Context) log.Fielder
type FieldProvider func() log.Fielder

type RequestLog struct {
	providers []RequestFieldProvider
	entry     log.Interface
	ctn       di.Container
}

func New(ctn di.Container, provider ...RequestFieldProvider) *RequestLog {
	return &RequestLog{providers: provider, ctn: ctn}
}

func (l *RequestLog) Log(ctx *gin.Context) log.Interface {
	if l.entry == nil {
		entry := l.ctn.Get("log").(log.Interface)
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
