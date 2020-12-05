package spring

import (
	"github.com/sarulabs/di"

	"github.com/summer-solutions/orm"
)

type ORMRegistryInitFunc func(registry *orm.Registry)

func OrmRegistry(init ORMRegistryInitFunc) *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "orm_config",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService, err := ctn.SafeGet("config")
			if err != nil {
				return nil, err
			}
			registry := orm.InitByYaml(configService.(*ViperConfig).Get("orm").(map[string]interface{}))
			init(registry)
			ormConfig, err := registry.Validate()
			return ormConfig, err
		},
	}
}
