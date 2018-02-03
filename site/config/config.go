package config

import (
	"ems/core/redirect_back"
	"ems/mailer"
	"ems/render"
	"ems/session/manager"
	"html/template"
	"os"

	"github.com/jinzhu/configor"
	"github.com/microcosm-cc/bluemonday"
	"ems/mailer/logger"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Site     string
}

var Config = struct {
	Port uint `default:"5000" env:"PORT"`
	DB   struct {
		Name     string `env:"DBName" default:"ems"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
	SMTP SMTPConfig
}{}

var (
	Root         = os.Getenv("GOPATH") + "/src/ems/site"
	View         *render.Render
	Mailer       *mailer.Mailer
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

func init() {
	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml"); err != nil {
		panic(err)
	}

	View = render.New(nil)

	htmlSanitizer := bluemonday.UGCPolicy()
	View.RegisterFuncMap("raw", func(str string) template.HTML {
		return template.HTML(htmlSanitizer.Sanitize(str))
	})

	Mailer = mailer.New(&mailer.Config{
		Sender: logger.New(&logger.Config{}),
	})
}
