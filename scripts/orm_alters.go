package scripts

import (
	"os"

	"github.com/fatih/color"
	"github.com/summer-solutions/spring"
)

type ORMAltersScript struct {
}

func (script *ORMAltersScript) Active() bool {
	_, has := spring.DIC().OrmConfig()
	return has
}

func (script *ORMAltersScript) Unique() bool {
	return true
}

func (script *ORMAltersScript) Code() string {
	return "orm-alters"
}

func (script *ORMAltersScript) Description() string {
	return "orm-alters"
}

func (script *ORMAltersScript) Run() error {
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
		os.Exit(1)
	}
	return nil
}
