package services

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring/app"
	"github.com/summer-solutions/spring/ioc"
	"github.com/summer-solutions/spring/services/config"
)

func Config(configFolderPath string) *ioc.ServiceDefinition {
	return &ioc.ServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewViperConfig(ctn.Get("app").(*app.App).Name, configFolderPath)
		},
	}
}
