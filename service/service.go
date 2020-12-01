package service

import (
	"context"

	"github.com/summer-solutions/spring/service/log"

	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring/service/config"

	apexLog "github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
)

var container di.Container

func SetGlobalContainer(c di.Container) {
	container = c
}

func GetGlobalContainer() di.Container {
	return container
}

func GetRequestContainer(ctx context.Context) di.Container {
	c := ctx.Value("GinContextKey").(*gin.Context)

	requestContainer, has := c.Get("RequestContainer")
	if has {
		return requestContainer.(di.Container)
	}

	ioCRequestContainer, err := container.SubContainer()
	c.Set("RequestContainer", ioCRequestContainer)

	if err != nil {
		panic(err)
	}

	return ioCRequestContainer
}

func Log() apexLog.Interface {
	return GetGlobalContainer().Get("log").(apexLog.Interface)
}

func Config() *config.ViperConfig {
	return GetGlobalContainer().Get("config").(*config.ViperConfig)
}

func ConfigSafe() (*config.ViperConfig, bool) {
	v, err := GetGlobalContainer().SafeGet("config")
	return v.(*config.ViperConfig), err != nil
}

func OrmConfig() orm.ValidatedRegistry {
	return GetGlobalContainer().Get("orm_config").(orm.ValidatedRegistry)
}

func LogContext(ctx context.Context) *log.RequestLog {
	return GetRequestContainer(ctx).Get("log_request").(*log.RequestLog)
}

func OrmEngineContext(ctx context.Context) *orm.Engine {
	return GetRequestContainer(ctx).Get("orm_engine").(*orm.Engine)
}
