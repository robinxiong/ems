package auth

import (
	"ems/auth_themes/clean"
	"ems/auth"
	"ems/site/db"
	"ems/site/config"
	"ems/site/app/models"
)

var (
// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB: db.DB,
		Render: config.View,
		Mailer: config.Mailer,
		UserModel: models.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
	})
)

func init(){

}