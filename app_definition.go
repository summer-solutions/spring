package spring

const ModeLocal = "local"
const ModeProd = "prod"

type AppDefinition struct {
	mode string
	name string
}

func (app *AppDefinition) Name() string {
	return app.mode
}

func (app *AppDefinition) IsInLocalMode() bool {
	return app.IsInMode(ModeLocal)
}

func (app *AppDefinition) IsInProdMode() bool {
	return app.IsInMode(ModeProd)
}
func (app *AppDefinition) IsInMode(mode string) bool {
	return app.mode == mode
}
