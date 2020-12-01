package global

import (
	apexLog "github.com/apex/log"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/service/log"
)

func LogGlobalService(provider ...log.FieldProvider) spring.InitHandler {
	return func(s *spring.Server, def *spring.Def) {
		def.Name = "log"
		def.Build = func(ctn di.Container) (interface{}, error) {
			l := apexLog.WithFields(&apexLog.Fields{})
			for _, fields := range provider {
				l = l.WithFields(fields())
			}
			return l, nil
		}
	}
}
