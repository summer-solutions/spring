package spring

import (
	"context"

	ginLocal "github.com/summer-solutions/spring/gin"

	"github.com/summer-solutions/spring/services/log"

	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring/services/config"

	apexLog "github.com/apex/log"
	"github.com/sarulabs/di"
)

var container di.Container

func GetContainer() di.Container {
	return container
}

type DIServiceDefinition struct {
	Name   string
	Global bool
	Build  func(ctn di.Container) (interface{}, error)
	Close  func(obj interface{}) error
}

func GetContainerForRequest(ctx context.Context) di.Container {
	c := ginLocal.FromContext(ctx)

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

func DILog() (apexLog.Interface, bool) {
	v, err := GetContainer().SafeGet("log")
	if err == nil {
		return v.(apexLog.Interface), true
	}
	return nil, false
}

func DIConfig() (*config.ViperConfig, bool) {
	v, err := GetContainer().SafeGet("config")
	if err == nil {
		return v.(*config.ViperConfig), true
	}
	return nil, false
}

func DIOrmConfig() (orm.ValidatedRegistry, bool) {
	v, err := GetContainer().SafeGet("orm_config")
	if err == nil {
		return v.(orm.ValidatedRegistry), true
	}
	return nil, false
}

func DILogForContext(ctx context.Context) (*log.RequestLog, bool) {
	v, err := GetContainerForRequest(ctx).SafeGet("log_request")
	if err == nil {
		return v.(*log.RequestLog), true
	}
	return nil, false
}

func DIOrmEngineForContext(ctx context.Context) (*orm.Engine, bool) {
	v, err := GetContainerForRequest(ctx).SafeGet("orm_engine")
	if err == nil {
		return v.(*orm.Engine), true
	}
	return nil, false
}
