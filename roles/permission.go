package roles

import (
	"errors"
	"fmt"
)

//PermissionMode 权限模式字符串
type PermissionMode string

const (
	//Create 预定义权限模式, 表示创建权限
	Create PermissionMode = "create"
	//Read
	Read PermissionMode = "read"
	//Update
	Update PermissionMode = "update"
	//Delete
	Delete PermissionMode = "delete"
	//CRUD
	CRUD PermissionMode = "crud"
)

//ErrPermissionDenied on permission error
var ErrPermissionDenied = errors.New("permission denied")

//第一个参数中是否存在于第二个参数数组中
func includeRoles(roles []string, values []string) bool{
	for _, role := range roles {
		if role == Anyone {
			return true
		}
		for _, value := range values {
			if value == role {
				return true
			}
		}
	}
	return false
}

//Permission 主要用于一个资源的权限设置，哪些角色可以访问，哪些角色不可以访问
type Permission struct {
	Role *Role
	AllowedRoles map[PermissionMode][]string
	DeniedRoles map[PermissionMode][]string
}

//Allow allows permission mode for roles
func (permission *Permission) Allow(mode PermissionMode, roles ...string) *Permission{
	if mode == CRUD {
		return permission.Allow(Create, roles...).Allow(Update, roles...).Allow(Read, roles...).Allow(Delete, roles...)
	}

	if permission.AllowedRoles[mode] == nil {
		permission.AllowedRoles[mode] = []string{}
	}
	permission.AllowedRoles[mode] = append(permission.AllowedRoles[mode], roles...)
	return permission
}


func (permission *Permission) Deny(mode PermissionMode, roles ...string) *Permission{
	if mode == CRUD {

		return permission.Deny(Create, roles...).Deny(Update, roles...).Deny(Read, roles...).Deny(Delete, roles...)
	}
	if permission.DeniedRoles[mode] == nil {
		permission.DeniedRoles[mode] = []string{}
	}
	permission.DeniedRoles[mode] = append(permission.DeniedRoles[mode], roles...)
	return permission
}

func (permission Permission) HasPermission(mode PermissionMode, roles ...interface{}) bool {
	var roleNames []string
	for _, role := range roles {
		if r, ok := role.(string); ok {
			roleNames = append(roleNames, r)
		} else if roler, ok := role.(Roler); ok {
			roleNames = append(roleNames, roler.GetRoles()...)
		} else {
			fmt.Printf("invalid role %#v\n", role)
			return false
		}
	}


	if len(permission.DeniedRoles) != 0 {
		if DeniedRoles := permission.DeniedRoles[mode]; DeniedRoles != nil {
			if includeRoles(DeniedRoles, roleNames) {
				return false
			}
		}
	}

	//如果没有定义允许的角色，则直接返回true
	if len(permission.AllowedRoles) == 0 {
		return true
	}

	if AllowedRoles := permission.AllowedRoles[mode]; AllowedRoles != nil {
		if includeRoles(AllowedRoles, roleNames) {
			return true
		}
	}

	return false
}

//Concat 将两个权限合并为一个
//后者不覆盖前者
func (permission *Permission) Concat(newPermission *Permission) *Permission{
	var result = Permission{
		Role: Global,
		AllowedRoles: map[PermissionMode][]string{},
		DeniedRoles: map[PermissionMode][]string{},
	}

	var appendRoles = func(p *Permission){
		if p != nil {
			result.Role = p.Role

			for mode, roles := range p.DeniedRoles {
				result.DeniedRoles[mode] = append(result.DeniedRoles[mode], roles...)
			}

			for mode, roles := range p.AllowedRoles {
				result.AllowedRoles[mode] = append(result.AllowedRoles[mode], roles...)
			}
		}

	}

	appendRoles(newPermission)
	appendRoles(permission)

	return &result

}