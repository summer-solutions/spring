package services

import (
	"github.com/summer-solutions/spring/di"
	"github.com/summer-solutions/spring/services/config"
)

func Config(configFolderPath string) *di.ServiceDefinition {
	return &di.ServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func() (interface{}, error) {
			return config.NewViperConfig(di.App().Name, configFolderPath)
		},
	}
}
