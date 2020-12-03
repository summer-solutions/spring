package app

const ModeLocal = "local"
const ModeProd = "prod"

type App struct {
	Mode string
	Name string
}

func (app *App) IsInLocalMode() bool {
	return app.IsInMode(ModeLocal)
}

func (app *App) IsInProdMode() bool {
	return app.IsInMode(ModeProd)
}
func (app *App) IsInMode(mode string) bool {
	return app.Mode == mode
}
