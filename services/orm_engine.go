package services

import (
	"fmt"

	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring/ioc"
)

func OrmEngine() *ioc.ServiceDefinition {
	return &ioc.ServiceDefinition{
		Name:   "orm_engine",
		Global: false,
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
