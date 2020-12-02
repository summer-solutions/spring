package app

const ModeLocal = "local"
const ModeDev = "dev"
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

func (app *App) IsInDevMode() bool {
	return app.IsInMode(ModeDev)
}

func (app *App) IsInMode(mode string) bool {
	return app.Mode == mode
}

func (app *App) AppName(mode string) bool {
	return app.Mode == mode
}
