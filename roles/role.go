package roles

import "net/http"

const (
	//Anyone可以作为作何的角色，通用的一个角色
	Anyone = "*"
)

//Role 是一个struct 类型， 包含所有的角色的定义 definitions
type Role struct {
	definitions map[string]func(request *http.Request, user interface{}) bool
}

//New 用于初始化一个角色
func New() *Role{
	return &Role{}
}