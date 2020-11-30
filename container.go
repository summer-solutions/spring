package spring

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Container struct {
	ctx       context.Context
	providers map[interface{}]ServiceProvider
}

type ServiceProvider interface {
	Key() interface{}
	Get(ctx context.Context) interface{}
}

func (c *Container) New(ctx context.Context, provider ...ServiceProvider) *Container {
	instance := &Container{ctx: ctx, providers: make(map[interface{}]ServiceProvider)}
	for _, p := range provider {
		instance.RegisterService(p)
	}
	return instance
}

func (c *Container) Get(key interface{}) (val interface{}, has bool) {
	val = c.ctx.Value(key)
	if val == nil {
		provider, has := c.providers[key]
		if !has {
			return nil, false
		}
		val = provider.Get(c.ctx)
		c.ctx = context.WithValue(c.ctx, key, val)
	}
	return val, val != nil
}

func (c *Container) MustGet(key interface{}) interface{} {
	val, _ := c.Get(key)
	return val
}

func (c *Container) GetFromRequest(ctx *gin.Context, key string) (val interface{}, has bool) {
	val, has = ctx.Get(key)
	if !has {
		provider, has := c.providers[key]
		if !has {
			return nil, false
		}
		val = provider.Get(c.ctx)
		ctx.Set(key, val)
	}
	return val, has
}

func (c *Container) RegisterService(provider ServiceProvider) {
	c.providers[provider.Key()] = provider
}
