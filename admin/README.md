#route
Route提供了serveMux, 它提供了所有请求的处理入口
然后调用admin.Admin对像中注册的路由

```go
var Admin *admin.Admin
func init(){
	Admin = admin.New(&admin.AdminConfig{})

}
//将route.go中的serverMux添加到全局的路径下，调用serverMux.serverHTTP方法，匹配Admin对像中的router对像
admin.Admin.MountTo("/admin", mux)

serveMux包含admin,admin包含router
```

##serverMux serveHTTP
1. 如果当前的请求路径是以assets开头，则返回相应的静态文件

2. 检查当前用户信息，如果没有相应的用户，则转向到登陆页面
   
   在config/admin/admin.go中指定一个认证接口，它需要实现admin/auth.go中的接口, 具体的struct为config/auth/admin_auth.go
   
3. 调用middleware，执行匹配到的handler

#Auth
主要是core/auth和core/auth_themes, auth_themes用于区分不同的语言，而auth package实现认证的逻辑


