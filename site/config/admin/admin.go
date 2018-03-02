package admin

import (
	"ems/admin"
	"ems/site/db"
	"ems/publish2"
	"ems/site/config/auth"
)

var Admin *admin.Admin

func init(){
	Admin = admin.New(&admin.AdminConfig{Auth: auth.AdminAuth{}, DB: db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)})

}