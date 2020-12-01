package request

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/service/log"
)

func LogRequestService(provider ...log.RequestFieldProvider) spring.InitHandler {
	return func(s *spring.Server, def *spring.Def) {
		def.Name = "log_request"
		def.Build = func(ctn di.Container) (interface{}, error) {
			l := log.New(ctn, provider...)
			return l, nil
		}
	}
}
