package scripts

import (
	"context"

	"github.com/sarulabs/di"

	"github.com/fatih/color"
	"github.com/summer-solutions/spring"
)

func ORMAlters() *spring.ServiceDefinition {
	return &spring.ServiceDefinition{
		Name:   "orm-alters",
		Global: true,
		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ORMAltersScript{}, nil
		},
	}
}

type ORMAltersScript struct {
}

func (script *ORMAltersScript) Active() bool {
	_, has := spring.DIC().OrmConfig()
	return has
}

func (script *ORMAltersScript) Unique() bool {
	return true
}

func (script *ORMAltersScript) Description() string {
	return "show all MySQL schema changes"
}

func (script *ORMAltersScript) Run(_ context.Context, exit spring.Exit) {
	ormEngine, _ := spring.DIC().OrmEngine()
	alters := ormEngine.GetAlters()
	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
	}
	if len(alters) > 0 {
		exit.Error()
	}
}
