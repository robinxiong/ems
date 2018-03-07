package account

import (
	"ems/site/config/application"
	"ems/site/config/auth"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}
// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
}


func (app App) ConfigureApplication(application *application.Application) {
	application.Router.Mount("/auth/", auth.Auth.NewServeMux())
}