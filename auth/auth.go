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
	"ems/session/manager"
	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	*Config
	providers []Provider  //登录的方式，在点击登录时，需要获取到auth中注册的provider, 即auth.GetProvider

	//嵌入SessionStorer，使它实现Authority's的AuthInterface接口, 一个Get方法，返回Claims
	//在new方法中，它的值为config的SessionStorer
	SessionStorerInterface
}
type Config struct {
	DB *gorm.DB
	//将URLPrefix挂载到router中去, 默认为/auth
	URLPrefix string
	// AuthIdentityModel 是一个数据库表，用来保存认证信息，比如email/password, OAuth token, 用户ID
	// 同时记录登录时间，IP， signin次数
	AuthIdentityModel interface{}
	//UserModel是一个user model struct, 它保存用户信息，它可以为空，为空时，Auth将假设认证信息没有用户
	//UserModel在config/auth/auth.go中指定为users.User{}
	UserModel  interface{}
	//UserStorer是一个接口,定义了如何获取/保存一个用户, 它将在provider/password/DefaultAuthorizeHandler中使用
	UserStorer UserStorerInterface

	//使用Render（https://github.com/qor/render）来读取登陆页面
	Render *render.Render
	//使用Mailer来发送认证的邮件, 它需要传递给provider
	Mailer *mailer.Mailer

	Redirector RedirectorInterface

	//添加额外搜索auth相关模板的路径. 比如login.tmpl，它会对过config.Render.RegisterViewPath向模板擎注册路径
	//否则模板引擎只在默认的app/views或者gopath, 以及app下的vendor中查找
	ViewPaths []string

	//auth/password/login 验证帐号密码，通常在provider/password/password.go Login方法中调用
	//第一个参数为auth/context它包含了auth, request, response, provider, claims
	//第二个参数为provider的验证函数
	//通常LoginHandler在auth.New中初始化，使用默认的DefaultLoginHandler
	LoginHandler func(*Context, func(*Context) (*claims.Claims, error))
	RegisterHandler func(*Context, func(*Context)(*claims.Claims, error))

	//SessionStorer 是一个接口，定义了对sesssion数据的编码，校验，保存，删除等，同时flash message
	//Auth提供了一个默认的方法来做这件事件, 在使用它之前，需要将SessionManager中件间添加到router中
	//更多可以参考https://github.com/qor/session
	SessionStorer SessionStorerInterface
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

	//获取和保存用户的接口
	if config.UserStorer == nil {
		config.UserStorer = &UserStorer{}
	}

	if config.SessionStorer == nil {
		config.SessionStorer = &SessionStorer{
			SessionName: "_auth_session", //保存认证数据的session名字, 它是_session下面的_auth_session
			SessionManager:manager.SessionManager,
			SigningMethod:  jwt.SigningMethodHS256,
		}
	}

	//处理登录的Handler
	if config.LoginHandler == nil {
		config.LoginHandler = DefaultLoginHandler
	}
	if config.RegisterHandler == nil {
		config.RegisterHandler = DefaultRegisterHandler
	}

	//先查找auth_themes中的模板，因粉auth_themes的New方法，添加了路径到config.ViewPath
	for _, viewPath := range config.ViewPaths {
		config.Render.RegisterViewPath(viewPath)
	}

	config.Render.RegisterViewPath("../auth/views")

	auth := &Auth{Config: config}

	return auth
}
