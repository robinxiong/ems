package admin

import "ems/site/config/application"

func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/admin"
	}
	return &App{Config: config}
}

type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
	Prefix string
}


func (app App) ConfigureApplication(application *application.Application) {
	Admin := application.Admin
	application.Router.Mount(app.Config.Prefix, Admin.NewServeMux(app.Config.Prefix))
}