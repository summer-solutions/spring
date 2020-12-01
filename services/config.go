package services

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/services/config"
)

func Config(configFilePath string) *spring.DIServiceDefinition {
	return &spring.DIServiceDefinition{
		Name:   "config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewViperConfig(configFilePath)
		},
	}
}
