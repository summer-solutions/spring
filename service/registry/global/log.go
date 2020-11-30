package global

import (
	"github.com/apex/log"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
)

var LogGlobalService spring.InitHandler = func(s *spring.Server, def *spring.Def) {
	def.Name = "log"
	def.Build = func(ctn di.Container) (interface{}, error) {
		return log.WithFields(&log.Fields{}), nil
	}
}
