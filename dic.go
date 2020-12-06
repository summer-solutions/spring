package spring

import (
	"context"

	apexLog "github.com/apex/log"
	"github.com/summer-solutions/orm"
)

type DICInterface interface {
	App() *AppDefinition
	Log() apexLog.Interface
	Config() *ViperConfig
	OrmConfig() (orm.ValidatedRegistry, bool)
	LogForContext(ctx context.Context) *RequestLog
	OrmEngineForContext(ctx context.Context) (*orm.Engine, bool)
}

type dic struct {
}

var dicInstance *dic

func DIC() DICInterface {
	return dicInstance
}

func (d *dic) App() *AppDefinition {
	return GetServiceRequired("app").(*AppDefinition)
}

func (d *dic) Log() apexLog.Interface {
	return GetServiceRequired("log").(apexLog.Interface)
}

func (d *dic) Config() *ViperConfig {
	return GetServiceRequired("config").(*ViperConfig)
}

func (d *dic) OrmConfig() (orm.ValidatedRegistry, bool) {
	v, has := GetServiceOptional("orm_config")
	if has {
		return v.(orm.ValidatedRegistry), true
	}
	return nil, false
}

func (d *dic) LogForContext(ctx context.Context) *RequestLog {
	return GetServiceForRequestRequired(ctx, "log_request").(*RequestLog)
}

func (d *dic) OrmEngineForContext(ctx context.Context) (*orm.Engine, bool) {
	v, has := GetServiceForRequestOptional(ctx, "orm_engine")
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}