package clean

import (
	"ems/auth"
	"ems/auth/providers/password"
	"ems/auth/claims"
	"errors"
)
var ErrPasswordConfirmationNotMatch = errors.New("password confirmation doesn't match password")
func New(config *auth.Config) *auth.Auth{
	//读取auth_themes下的模板文件，否则查找auth/views下的文件
	if config == nil {
		config = &auth.Config{}
	}
	config.ViewPaths = append(config.ViewPaths, "../auth_themes/clean/views")
	Auth := auth.New(config)


	//帐号登录和注册的功能
	Auth.RegisterProvider(password.New(&password.Config{
		Confirmable: false,
		RegisterHandler: func(context *auth.Context) (*claims.Claims, error) {
			context.Request.ParseForm()

			if context.Request.Form.Get("confirm_password") != context.Request.Form.Get("password") {
				return nil, ErrPasswordConfirmationNotMatch
			}

			return password.DefaultRegisterHandler(context)
		},
	}))
	return Auth
}