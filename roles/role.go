package roles

import (
	"net/http"

	"fmt"
)

const (
	//Anyone可以作为作何的角色，通用的一个角色
	Anyone = "*"
)

//Checker 检查当前的request是否跟跟当前的角色匹配
type Checker func(req *http.Request, user interface{}) bool

//Role 是一个struct 类型， 包含所有的角色的定义 definitions
type Role struct {
	definitions map[string]Checker
}

//New 用于初始化一个角色
func New() *Role{
	return &Role{}
}


//Register register role with confitions
func (role *Role) Register(name string, fc Checker){
	if role.definitions == nil {
		role.definitions = map[string]Checker{}
	}
	definition :=role.definitions[name]

	if definition != nil {
		fmt.Printf("%v already defined, overwrited it!\n", name)
	}

	role.definitions[name] = fc
}


//NewPermission 初妈化角色的权限
func (role *Role) NewPermission() *Permission{
	return &Permission{
		Role: role,
		AllowedRoles:map[PermissionMode][]string{},
		DeniedRoles:map[PermissionMode][]string{},
	}
}

//Allow 为角色添加允许权限模式

func (role *Role) Allow(mode PermissionMode, roles ...string) *Permission{
	return role.NewPermission()
}