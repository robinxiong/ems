package application

import (
	"github.com/jinzhu/gorm"
	"github.com/go-chi/chi"
	"net/http"
	"ems/assetfs"
	"ems/admin"
	"ems/middlewares"
	"ems/wildcard_router"
)
type MicroAppInterface interface {
	ConfigureApplication(*Application)
}
type Application struct {
	*Config
}

type Config struct {
	Router   *chi.Mux
	Handlers []http.Handler
	AssetFS  assetfs.Interface
	Admin    *admin.Admin
	DB       *gorm.DB
}

// New new application
func New(cfg *Config) *Application {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Router == nil {
		cfg.Router = chi.NewRouter()
	}

	if cfg.AssetFS == nil {
		cfg.AssetFS = assetfs.AssetFS()
	}

	return &Application{
		Config: cfg,
	}
}


// Use mount router into micro app
func (application *Application) Use(app MicroAppInterface) {
	app.ConfigureApplication(application)
}

// NewServeMux allocates and returns a new ServeMux.
func (application *Application) NewServeMux() http.Handler {
	if len(application.Config.Handlers) == 0 {
		return middlewares.Apply(application.Config.Router)
	}

	wildcardRouter := wildcard_router.New()
	for _, handler := range application.Config.Handlers {
		wildcardRouter.AddHandler(handler)
	}
	wildcardRouter.AddHandler(application.Config.Router)

	return middlewares.Apply(wildcardRouter)
}