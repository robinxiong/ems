package admin

import "ems/admin"

func SetupDashboard(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin", Priority: 1})
}