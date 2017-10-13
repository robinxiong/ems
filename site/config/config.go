package config

import (
	"os"
	"ems/render"
)

type SMTPConfig struct {
	Host string
	Port string
	User string
	Password string
	Site string
}

var Config = struct {
	Port uint `default:"7000" env:"PORT"`
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
	Root = os.Getenv("GOPATH") + "/src/ems/site"
	View *render.Render

)