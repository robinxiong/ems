package admin

import (
	"ems/core"
	"ems/roles"
	"path"
)

//GetMenus返回admin的所有左侧菜单
func (admin Admin) GetMenus() []*Menu {

	return admin.menus
}

//Usage: site/app/admin/dashboard.go Admin.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin", Priority: 1})
func (admin *Admin) AddMenu(menu *Menu) *Menu {
	menu.router = admin.router
	names := append(menu.Ancestors, menu.Name) //添加到完整名称中, 但并没有添加到Ancestors, 所以menu.Ancestors依然为0

	//如果菜单已经存在于admin.Menus中
	if old := admin.GetMenu(names...); old != nil {
		//当前菜单是多级菜单(一级菜单,names为1) 或者存在的菜单Ancestors为0（一级菜单), 替换原来的menu
		//如原来的菜单为手动添加的二级菜单（old.Ancestors为0），现在添加一个相同的一级菜单时(len(names) ==1)，则不替换
		//可以查看menu_test.go menu2
		if len(names) > 1 || len(old.Ancestors) == 0 {

			old.Link = menu.Link
			old.Priority = menu.Priority
			old.RelativePath = menu.RelativePath
			old.Permissioner = menu.Permissioner
			old.Permission = menu.Permission
			*menu = *old
			return old
		}
	}
	admin.menus = appendMenu(admin.menus, menu.Ancestors, menu)
	return menu
}

// GetMenu 通过name获取左侧的菜单, name可以是单个的数组，也可以是包含父菜单的名称的数组
func (admin Admin) GetMenu(name ...string) *Menu {
	return getMenu(admin.menus, name...)
}

/*
	Menu admin界面左侧菜单定义
	Name用来标识一个菜单，而Ancestors用来表示当前菜单的路径Product Managements/Products, subMenus是一个menu数组，用来保存它的子菜单
	当要查找一个菜单时，比如从 []menus{menu1, menu2, menu3}中查找Product Managements/Products, 首先它在一级菜单中查找Product Managements
	然后在Product Managements的子菜单subMenus中找到Products.
*/
type Menu struct {
	Name         string   //显示的名称
	Link         string   //菜单链接, 决对路径
	Priority     int      //优先级， 可以是负数，0，和正数，正数的优先级最高，其次是0，然后是负数，正数越小越靠前 1 > 5，负数也是越小越靠前 -5 > -1, 如果都默认为0，则先添加，先优先
	router       *Router  //绑定admin.Router到当前Menu
	Ancestors    []string //手动指定Menu的上级目录, 它主要用于resource， 可以参考admin.AddResource方法，为resource添加到menu
	subMenus     []*Menu  //子菜单
	RelativePath string   //相对路径，加上router.Prefix组合成决对路径

	Permissioner HasPermissioner   //检查权限的接口
	Permission   *roles.Permission //当前Menu的权限
}

/*
	URL 返回菜单的URL
	如果指定了Link，则直接返回, 如果没有指定Link, 则使用 admin.Router.Prefix + Menu的RelativePath
*/
func (menu *Menu) URL() string {
	if menu.Link != "" {
		return menu.Link
	}
	if (menu.router != nil) && (menu.RelativePath != "") {
		return path.Join(menu.router.Prefix, menu.RelativePath)
	}

	return menu.RelativePath
}

// GetSubMenus get submenus for a menu
func (menu *Menu) GetSubMenus() []*Menu {
	return menu.subMenus
}

func (menu *Menu) HasPermission(mode roles.PermissionMode, context *core.Context) bool {
	if menu.Permission != nil {
		var roles = []interface{}{}
		for _, role := range context.Roles {
			roles = append(roles, role)
		}
		return menu.Permission.HasPermission(mode, roles...)
	}

	if menu.Permissioner != nil {
		return menu.Permissioner.HasPermission(mode, context)
	}

	return true
}

//从Menu数据中，根据names来找到menu, names是有顺序的[中国，广东，深圳]
func getMenu(menus []*Menu, names ...string) *Menu {
	if len(names) > 0 {
		name := names[0]
		for _, menu := range menus {
			//还不是当前要查找的menu, 而是当前menu的父menu
			if len(names) > 1 {

			} else {
				//等于1时
				//找到当前name的menu
				if menu.Name == name {
					return menu
				}
				//只输入了 深圳，则遍历所有的menu, 以及它的子菜单，如果找到，则返回menu
				if len(menu.subMenus) > 0 {
					if m := getMenu(menu.subMenus, name); m != nil {
						return m
					}
				}
			}
		}
	}
	return nil
}

//appendMenu中， 传入的ancestors是不同的，【中国，广东，深圳]
//如果之前没有创建过它的父菜单，则先创建直接父目录，在一级级往上创建. 根据ancestors的名称，倒序创建
func generateMenu(ancestors []string, menu *Menu) *Menu {
	menuCount := len(ancestors)
	for index := range ancestors {
		menu = &Menu{Name: ancestors[menuCount-index-1], subMenus: []*Menu{menu}}
	}

	return menu
}

//将menu添加到Menu中指定的位置，位置由ancestors确定【中国，广东], 深圳
func appendMenu(menus []*Menu, ancestors []string, menu *Menu) []*Menu {
	//如果有指定位置
	if len(ancestors) > 0 {
		for _, m := range menus {
			if m.Name != ancestors[0] {
				continue
			}

			//将当前菜单，添加到m.subMenus中, 并返回整个mneus
			//如果指定的目录不是一级目录, 则添加到后面的菜单
			//否则，创建一个空的slice, 就不会在递归调用appendMenu方法
			if len(ancestors) > 1 {
				m.subMenus = appendMenu(m.subMenus, ancestors[1:], menu)
			} else {
				m.subMenus = appendMenu(m.subMenus, []string{}, menu)
			}

			return menus
		}
	}

	//ancestors大于0，而且没有找到父目录的情况, 或者ancestors为空，否则直接执行到上一步, 将menu添加到指定位置后，返回
	var newMenu = generateMenu(ancestors, menu)
	var added bool

	//如果当前不存在其它menus，则直接添加
	if len(menus) == 0 {
		menus = append(menus, newMenu)
	} else if newMenu.Priority > 0 {
		//新的Menu带了Priority属性， 并且大于0
		for idx, menu := range menus {
			//找到大于当前的menu，或者负数，或者0的menu，后移到当前新的menu后
			if menu.Priority > newMenu.Priority || menu.Priority <= 0 {
				menus = append(menus[0:idx], append([]*Menu{newMenu}, menus[idx:]...)...)
				added = true
				break
			}
		}
		//没有找到，则直接添加
		if !added {
			menus = append(menus, menu)
		}
	} else if newMenu.Priority < 0 {
		//从原来排序好的menus末尾开始查找
		for idx := len(menus) - 1; idx >= 0; idx-- {
			menu := menus[idx]
			if menu.Priority < newMenu.Priority || menu.Priority == 0 {
				menus = append(menus[0:idx+1], append([]*Menu{newMenu}, menus[idx+1:]...)...)
				added = true
				break
			}
		}

		if !added {
			menus = append(menus, menu)
		}
	} else {

		//如果默认是0的，则正数和原来为0的menu后
		for idx := len(menus) - 1; idx >= 0; idx-- {
			menu := menus[idx]
			if menu.Priority >= 0 {
				menus = append(menus[0:idx+1], append([]*Menu{newMenu}, menus[idx+1:]...)...)
				added = true
				break
			}
		}

		if !added {
			menus = append([]*Menu{menu}, menus...)
		}
	}
	return menus
}
