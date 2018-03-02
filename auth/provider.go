package auth

import "fmt"


//Provider提供了自己的路由，用来处理登入，注册， 比如auth/password/login就是以密码来登录系统
type Provider interface {
	GetName() string
	ConfigAuth(*Auth)
	Login(*Context)
	Logout(*Context)
	Register(*Context)
	Callback(*Context)
	ServeHTTP(*Context)
}


func (auth *Auth) RegisterProvider(provider Provider){
	name := provider.GetName()
	//查找是否已经注册
	for _, p := range auth.providers {
		if p.GetName() == name {
			fmt.Printf("warning: auth provider %v already registered", name)
			return
		}
	}
	provider.ConfigAuth(auth)
	auth.providers = append(auth.providers, provider)
}

// GetProvider get provider with name
func (auth *Auth) GetProvider(name string) Provider {
	for _, provider := range auth.providers {
		if provider.GetName() == name {
			return provider
		}
	}
	return nil
}

// GetProviders return registered providers
func (auth *Auth) GetProviders() (providers []Provider) {
	for _, provider := range auth.providers {
		providers = append(providers, provider)
	}
	return
}