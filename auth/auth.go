package auth

import (
	"ems/auth/auth_identity"
	"ems/mailer"
	"ems/mailer/logger"
	"ems/render"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"ems/auth/claims"
)

type Auth struct {
	*Config
	providers []Provider  //登录的方式，在点击登录时，需要获取到auth中注册的provider, 即auth.GetProvider
	//auth/password/login 验证帐号密码，通常在provider/password/password.go Login方法中调用
	//第一个参数为auth/context它包含了auth, request, response, provider, claims
	//第二个参数为provider的验证函数
	//通常LoginHandler在auth.New中初始化，使用默认的DefaultLoginHandler
	LoginHandler func(*Context, func(*Context) (*claims.Claims, error))
	//SessionStorer 是一个接口，定义了对sesssion数据的编码，校验，保存，删除等，同时flash message
	//Auth提供了一个默认的方法来做这件事件, 在使用它之前，需要将SessionManager中件间添加到router中
	//更多可以参考https://github.com/qor/session

}
type Config struct {
	DB *gorm.DB
	//将URLPrefix挂载到router中去, 默认为/auth
	URLPrefix string
	// AuthIdentityModel 是一个数据库表，用来保存认证信息，比如email/password, OAuth token, 用户ID
	// 同时记录登录时间，IP， signin次数
	AuthIdentityModel interface{}

	//使用Render（https://github.com/qor/render）来读取登陆页面
	Render *render.Render
	//使用Mailer来发送认证的邮件, 它需要传递给provider
	Mailer *mailer.Mailer
	//UserStorer用来提供get/save user, 默认提供基于AuthIdentityModel和UserModel两个struct
	UserModel  interface{}
	Redirector RedirectorInterface

	//添加额外搜索auth相关模板的路径. 比如login.tmpl，它会对过config.Render.RegisterViewPath向模板擎注册路径
	//否则模板引擎只在默认的app/views或者gopath, 以及app下的vendor中查找
	ViewPaths []string
}

//New 初始化 Auth
/*
	admin_auth
	// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB:         db.DB,
		Render:     config.View,
		Mailer:     config.Mailer,
		UserModel:  models.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: Auth,
	})
*/
func New(config *Config) *Auth {
	if config == nil {
		config = &Config{}
	}

	if config.URLPrefix == "" {
		config.URLPrefix = "/auth/"
	} else {
		config.URLPrefix = fmt.Sprintf("/%v/", strings.Trim(config.URLPrefix, "/"))
	}
	if config.AuthIdentityModel == nil {
		config.AuthIdentityModel = &auth_identity.AuthIdentity{}
	}

	if config.Render == nil {
		config.Render = render.New(nil)
	}

	if config.Mailer == nil {
		config.Mailer = mailer.New(&mailer.Config{
			Sender: logger.New(&logger.Config{}),
		})
	}

	//先查找auth_themes中的模板，因粉auth_themes的New方法，添加了路径到config.ViewPath
	for _, viewPath := range config.ViewPaths {
		config.Render.RegisterViewPath(viewPath)
	}

	config.Render.RegisterViewPath("../auth/views")

	auth := &Auth{Config: config}
	return auth
}
