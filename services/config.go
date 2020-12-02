package services

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring/app"
	diLocal "github.com/summer-solutions/spring/di"
	"github.com/summer-solutions/spring/services/config"
)

func Config(configFolderPath string) *diLocal.ServiceDefinition {
	return &diLocal.ServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewViperConfig(ctn.Get("app").(*app.App).Name, configFolderPath)
		},
	}
}
