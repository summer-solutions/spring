package spring

import (
	"flag"
	"os"

	"github.com/fatih/color"
)

func PreDeploy() {
	preDeployFlag := flag.Bool("pre-deploy", false, "Execute pre deploy mode")
	flag.Parse()

	if !*preDeployFlag {
		return
	}

	ormConfigService, has := DIC().OrmConfig()
	if !has {
		return
	}
	ormService := ormConfigService.CreateEngine()

	alters := ormService.GetAlters()

	hasAlters := false
	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
		hasAlters = true
	}

	if hasAlters {
		os.Exit(1)
	}

	if !DIC().App().IsInLocalMode() {
		os.Exit(0)
	}
}
