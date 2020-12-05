package services

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/services/config"
)

func Config(configFolderPath string) *spring.ServiceDefinition {
	return &spring.ServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewViperConfig(ctn.Get("app").(*spring.AppDefinition).Name, configFolderPath)
		},
	}
}
