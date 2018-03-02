package password

import (
	"ems/auth"
	"ems/auth/providers/password/encryptor"
	"ems/auth/providers/password/encryptor/bcrypt_encryptor"
)

type Config struct {
	Encryptor encryptor.Interface //加密密码和密码比对的接口
}


func New(config *Config) *Provider {
	if config == nil {
		config = &Config{}
	}

	//如果没有指定加密方法，则使用encryptor/bcrypt_encryptor来加密
	if config.Encryptor == nil {
		config.Encryptor = bcrypt_encryptor.New(&bcrypt_encryptor.Config{})
	}

	provider := &Provider{Config: config}



	return provider
}


type Provider struct {
	*Config
}

func (*Provider) GetName() string {
	return "password"
}
//auth AddProvider方法调用
func (provider *Provider) ConfigAuth(auth *auth.Auth) {
	auth.Render.RegisterViewPath("/auth/providers/password/views")
	if auth.Mailer != nil {
		auth.Mailer.RegisterViewPath("/auth/providers/password/views/mailers")
	}
}

func (*Provider) Login(context *auth.Context) {
	context.Auth.LoginHandler()
}

func (*Provider) Logout(*auth.Context) {
	panic("implement me")
}

func (*Provider) Register(*auth.Context) {
	panic("implement me")
}

func (*Provider) Callback(*auth.Context) {
	panic("implement me")
}

func (*Provider) ServeHTTP(*auth.Context) {
	panic("implement me")
}