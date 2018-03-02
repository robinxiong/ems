package clean

import (
	"ems/auth"
	"ems/auth/providers/password"
)

func New(config *auth.Config) *auth.Auth{
	//读取auth_themes下的模板文件，否则查找auth/views下的文件
	if config == nil {
		config = &auth.Config{}
	}
	config.ViewPaths = append(config.ViewPaths, "../auth_themes/clean/views")
	Auth := auth.New(config)


	//帐号登录和注册的功能
	Auth.RegisterProvider(password.New(&password.Config{

	}))
	return Auth
}