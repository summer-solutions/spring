package request

import (
	"fmt"

	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/service"

	"github.com/sarulabs/di"
)

var OrmEngineRequestService spring.InitHandler = func(s *spring.Server, def *spring.Def) {

	def.Name = "orm_engine"
	def.Build = func(ctn di.Container) (interface{}, error) {
		ormConfigService, has := service.OrmConfig()
		if !has {
			return nil, fmt.Errorf("missing orm config service")
		}
		ormEngine := ormConfigService.CreateEngine()
		ormEngine.SetLogMetaData("Source", "web-api")
		return ormEngine, nil
	}
}
