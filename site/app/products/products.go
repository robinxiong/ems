package products

import (
	"ems/admin"
	"ems/site/config/application"
	"ems/site/models/products"
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
	app.ConfigureAdmin(application.Admin)
}

func (App) ConfigureAdmin(Admin *admin.Admin) {

	Admin.AddMenu(&admin.Menu{Name: "Product Management", Priority: 1})

	Admin.AddResource(&products.Color{}, &admin.Config{Menu: []string{"Product Management"}, Priority: -5})
}
