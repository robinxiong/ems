package auth

import (
	"ems/admin"
	"ems/core"
	"ems/roles"
	"ems/site/app/models"
	"fmt"
	"net/http"
)

func init() {
	roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(*models.User).Role == "Admin"
	})
}

type AdminAuth struct {
}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) core.CurrentUser {
	currentUser := Auth.GetCurrentUser(c.Request)
	if currentUser != nil {
		coreCurrentUser, ok := currentUser.(core.CurrentUser)
		if !ok {
			fmt.Printf("User %#v haven't implement core.CurrentUser interface\n", currentUser)
		}
		return coreCurrentUser
	}

	return nil
}
