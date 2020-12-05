package spring

const ModeLocal = "local"
const ModeProd = "prod"

type AppDefinition struct {
	Mode string
	Name string
}

func (app *AppDefinition) IsInLocalMode() bool {
	return app.IsInMode(ModeLocal)
}

func (app *AppDefinition) IsInProdMode() bool {
	return app.IsInMode(ModeProd)
}
func (app *AppDefinition) IsInMode(mode string) bool {
	return app.Mode == mode
}
