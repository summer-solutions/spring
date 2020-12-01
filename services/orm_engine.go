package services

import (
	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring"

	"github.com/sarulabs/di"
)

func OrmEngine() *spring.CIServiceDefinition {
	return &spring.CIServiceDefinition{
		Name:   "orm_engine",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfigService, err := ctn.SafeGet("orm_config")
			if err != nil {
				return nil, err
			}
			ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
			return ormEngine, nil
		},
	}
}
