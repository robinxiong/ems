package password

import (
	"ems/auth"
	"ems/auth/providers/password/encryptor"
	"ems/auth/providers/password/encryptor/bcrypt_encryptor"
	"ems/auth/claims"
)

type Config struct {
	Confirmable bool //provider是否开启验证，如果开启，则需要判断帐号是否已经验证过了，它在auth_themes/clean中注册时指定，在hadnler/DefaultAuthorizeHandler中使用
	ConfirmMailer  func(email string, context *auth.Context, claims *claims.Claims, currentUser interface{}) error
	Encryptor encryptor.Interface //加密密码和密码比对的接口
	AuthorizeHandler func(*auth.Context)(*claims.Claims, error) //验证用户名和密码， defaultAuthorizeHandler
	RegisterHandler  func(*auth.Context) (*claims.Claims, error)
}


func New(config *Config) *Provider {
	if config == nil {
		config = &Config{}
	}

	//如果没有指定加密方法，则使用encryptor/bcrypt_encryptor来加密
	if config.Encryptor == nil {
		config.Encryptor = bcrypt_encryptor.New(&bcrypt_encryptor.Config{})
	}

	if config.ConfirmMailer == nil {
		config.ConfirmMailer = DefaultConfirmationMailer
	}

	if config.AuthorizeHandler == nil {
		config.AuthorizeHandler = DefaultAuthorizeHandler
	}
	if config.RegisterHandler == nil {
		config.RegisterHandler = DefaultRegisterHandler
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

func (provider *Provider) Login(context *auth.Context) {
	context.Auth.LoginHandler(context, provider.AuthorizeHandler)
}

func (*Provider) Logout(*auth.Context) {
	panic("implement me")
}

func (provider *Provider) Register(context *auth.Context) {
	context.Auth.RegisterHandler(context, provider.RegisterHandler)
}

func (*Provider) Callback(*auth.Context) {
	panic("implement me")
}

func (*Provider) ServeHTTP(*auth.Context) {
	panic("implement me")
}