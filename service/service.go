package service

import (
	"context"

	"github.com/summer-solutions/spring/service/config"

	"github.com/apex/log"
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

func Log() log.Interface {
	return GetGlobalContainer().Get("log").(log.Interface)
}

func Config() *config.ViperConfig {
	return GetGlobalContainer().Get("config").(*config.ViperConfig)
}
