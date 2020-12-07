package spring

import (
	"fmt"

	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

func ServiceDefinitionOrmEngine() []*ServiceDefinition {
	result := make([]*ServiceDefinition, 2)
	i := 0
	for b, name := range map[bool]string{true: "global", false: "request"} {
		def := &ServiceDefinition{
			Name:   "orm_engine_" + name,
			Global: b,
			Build: func(ctn di.Container) (interface{}, error) {
				ormConfigService, err := ctn.SafeGet("orm_config")
				if err != nil {
					return nil, fmt.Errorf("missing orm config service")
				}
				ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
				return ormEngine, nil
			},
		}
		result[i] = def
		i++
	}
	return result
}
