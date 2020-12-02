package services

import (
	"fmt"

	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring/di"
)

func OrmEngine() *di.ServiceDefinition {
	return &di.ServiceDefinition{
		Name:   "orm_engine",
		Global: false,
		Build: func() (interface{}, error) {
			ormConfigService, has := di.OrmConfig()
			if !has {
				return nil, fmt.Errorf("missing orm config service")
			}
			ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
			return ormEngine, nil
		},
	}
}
