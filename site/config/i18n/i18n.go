package i18n

import (
	"ems/i18n"
	"ems/i18n/backends/database"
	"ems/site/db"
	"ems/i18n/backends/yaml"
	"path/filepath"
	"ems/site/config"
)

var I18n *i18n.I18n

func init(){
	I18n = i18n.New(database.New(db.DB), yaml.New(filepath.Join(config.Root, "config/locales")))
}