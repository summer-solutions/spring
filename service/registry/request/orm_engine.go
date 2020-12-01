package request

import (
	"github.com/summer-solutions/orm"

	"github.com/summer-solutions/spring"

	"github.com/sarulabs/di"
)

var OrmEngineRequestService spring.InitHandler = func(s *spring.Server, def *spring.Def) {

	def.Name = "orm_engine"
	def.Build = func(ctn di.Container) (interface{}, error) {
		ormConfigService, err := ctn.SafeGet("orm_config")
		if err != nil {
			return nil, err
		}
		ormEngine := ormConfigService.(orm.ValidatedRegistry).CreateEngine()
		return ormEngine, nil
	}
}
