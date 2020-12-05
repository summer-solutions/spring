package services

import (
	"github.com/sarulabs/di"
	"github.com/summer-solutions/spring"

	"github.com/summer-solutions/orm"
)

type RegistryInitFunc func(registry *orm.Registry)

func OrmRegistry(init RegistryInitFunc) *spring.ServiceDefinition {
	return &spring.ServiceDefinition{
		Name:   "orm_config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService, err := ctn.SafeGet("config")
			if err != nil {
				return nil, err
			}
			registry := orm.InitByYaml(configService.(*spring.ViperConfig).Get("orm").(map[string]interface{}))
			init(registry)
			ormConfig, err := registry.Validate()
			return ormConfig, err
		},
	}
}
