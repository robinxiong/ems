package roles

import "github.com/pkg/errors"

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

//Permission 主要用于一个资源的权限设置，哪些角色可以访问，哪些角色不可以访问
type Permission struct {

}

func NewPermission() *Permission