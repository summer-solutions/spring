package spring

import (
	"fmt"

	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

func ServiceDefinitionOrmEngine() *ServiceDefinition {
	return serviceDefinitionOrmEngine(true)
}

func ServiceDefinitionOrmEngineForContext() *ServiceDefinition {
	return serviceDefinitionOrmEngine(false)
}

func serviceDefinitionOrmEngine(global bool) *ServiceDefinition {
	suffix := "request"
	if global {
		suffix = "request"
	}
	return &ServiceDefinition{
		Name:   "orm_engine_" + suffix,
		Global: global,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfigService, err := ctn.SafeGet("orm_config")
			if err != nil {
				return nil, fmt.Errorf("missing orm config service")
			}
			ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
			return ormEngine, nil
		},
	}
}
