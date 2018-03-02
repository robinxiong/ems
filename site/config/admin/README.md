##Route
创建一个admin.Admin代对像
```go
var Admin *admin.Admin
func init(){
	Admin = admin.New(&admin.AdminConfig{})

}
admin.Admin.MountTo("/admin", mux)
```