package global

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/service/config"
)

var ConfigGlobalService spring.InitHandler = func(s *spring.Server, def *spring.Def) {
	def.Name = "config"
	def.Build = func(ctn di.Container) (interface{}, error) {
		return config.NewViperConfig("../../config/web-api/config.local.yaml")
	}
}
