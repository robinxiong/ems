package roles

import (
	"testing"
	"github.com/qor/roles"
	"net/http"

)

func TestAllow(t *testing.T){
	permission := roles.Allow(roles.Read, "api")
	if !permission.HasPermission(roles.Read, "api") {
		t.Errorf("API should has permission to Read")
	}

	if permission.HasPermission(roles.Update, "api") {
		t.Errorf("API should has no permission to Update")
	}

	if permission.HasPermission(roles.Read, "admin") {
		t.Errorf("admin should has no permission to Read")
	}

	if permission.HasPermission(roles.Update, "admin") {
		t.Errorf("admin should has no permission to Update")
	}
}

func TestDeny(t *testing.T) {
	permission := roles.Deny(roles.Create, "api")

	if !permission.HasPermission(roles.Read, "api") {
		t.Errorf("API should has permission to Read")
	}

	if !permission.HasPermission(roles.Update, "api") {
		t.Errorf("API should has permission to Update")
	}

	if permission.HasPermission(roles.Create, "api") {
		t.Errorf("API should has no permission to Update")
	}

	if !permission.HasPermission(roles.Read, "admin") {
		t.Errorf("admin should has permission to Read")
	}

	if !permission.HasPermission(roles.Create, "admin") {
		t.Errorf("admin should has permission to Update")
	}
}

func TestCRUD(t *testing.T) {
	permission := roles.Allow(roles.CRUD, "admin")
	if !permission.HasPermission(roles.Read, "admin") {
		t.Errorf("Admin should has permission to Read")
	}

	if !permission.HasPermission(roles.Update, "admin") {
		t.Errorf("Admin should has permission to Update")
	}

	if permission.HasPermission(roles.Read, "api") {
		t.Errorf("API should has no permission to Read")
	}

	if permission.HasPermission(roles.Update, "api") {
		t.Errorf("API should has no permission to Update")
	}
}

func TestAll(t *testing.T) {
	permission := roles.Allow(roles.Update, roles.Anyone)

	if permission.HasPermission(roles.Read, "api") {
		t.Errorf("API should has no permission to Read")
	}

	if !permission.HasPermission(roles.Update, "api") {
		t.Errorf("API should has permission to Update")
	}

	permission2 := roles.Deny(roles.Update, roles.Anyone)

	if !permission2.HasPermission(roles.Read, "api") {
		t.Errorf("API should has permission to Read")
	}

	if permission2.HasPermission(roles.Update, "api") {
		t.Errorf("API should has no permission to Update")
	}
}

func TestCustomizePermission(t *testing.T) {
	var customized roles.PermissionMode = "customized"
	permission := roles.Allow(customized, "admin")

	if !permission.HasPermission(customized, "admin") {
		t.Errorf("Admin should has customized permission")
	}

	if permission.HasPermission(roles.Read, "admin") {
		t.Errorf("Admin should has no permission to Read")
	}

	permission2 := roles.Deny(customized, "admin")

	if permission2.HasPermission(customized, "admin") {
		t.Errorf("Admin should has customized permission")
	}

	if !permission2.HasPermission(roles.Read, "admin") {
		t.Errorf("Admin should has no permission to Read")
	}
}

func TestRegisterRole(t *testing.T) {
	permission := roles.Allow(roles.Create, "admin")
	roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
		//可以能currentUser作更多的判断
		return req.RemoteAddr == "127.0.0.1" && currentUser == nil
	})
	httpRequest, err := http.NewRequest("get", "/", nil)
	if err != nil {
		t.Errorf("error:", err.Error())
	}
	httpRequest.RemoteAddr = "127.0.0.1"
	MatchedRoles := roles.MatchedRoles(httpRequest, nil)
	pass := permission.HasPermission(roles.Create, MatchedRoles...)
	if pass == false{
		t.Errorf("当前用户不是拥有本地角色")
	}
}