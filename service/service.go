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

func Log() (apexLog.Interface, bool) {
	v, err := GetGlobalContainer().SafeGet("log")
	if err == nil {
		return v.(apexLog.Interface), true
	}
	return nil, false
}

func Config() (*config.ViperConfig, bool) {
	v, err := GetGlobalContainer().SafeGet("config")
	if err == nil {
		return v.(*config.ViperConfig), true
	}
	return nil, false
}

func OrmConfig() (orm.ValidatedRegistry, bool) {
	v, err := GetGlobalContainer().SafeGet("orm_config")
	if err == nil {
		return v.(orm.ValidatedRegistry), true
	}
	return nil, false
}

func LogContext(ctx context.Context) (*log.RequestLog, bool) {
	v, err := GetGlobalContainer().SafeGet("log_request")
	if err == nil {
		return v.(*log.RequestLog), true
	}
	return nil, false
}

func OrmEngineContext(ctx context.Context) (*orm.Engine, bool) {
	v, err := GetGlobalContainer().SafeGet("orm_engine")
	if err == nil {
		return v.(*orm.Engine), true
	}
	return nil, false
}
