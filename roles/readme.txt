#roles
https://doc.getqor.com/plugins/roles.html

1.权限模式，有预定义5种read, write, update, delete, crud, 我也可以创建自己的模式customizer
通过global的allow, 我们允许某一个角色，访问上面权限，比如global.AllowedRoles['read']=[]string{"角色1", "角色2"}

2. 注册角色，在我们检查权限之前，需要确认当前用户的角色, role.go提供了helper方法Register方法
// Register roles based on some conditions
  roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
      return req.RemoteAddr == "127.0.0.1" || (currentUser.(*User) != nil && currentUser.(*User).Role == "admin")
  })

  roles.Register("user", func(req *http.Request, currentUser interface{}) bool {
    return currentUser.(*User) != nil
  })

  roles.Register("visitor", func(req *http.Request, currentUser interface{}) bool {
    return currentUser.(*User) == nil
  })

  // Get roles from a user
  matchedRoles := roles.MatchedRoles(httpRequest, user) // []string{"user", "admin"}

  // Check if role `user` or `admin` has Read permission
  permission.HasPermission(roles.Read, matchedRoles...)

  以上的代码都是先定义了一个Global的角色(global.go), 然后注册了三个角色admin, user, visitor, roles.MatchedRoles会调用register方法中的第二个回调函数，判断当前用户拥有哪些角色
  然后可以用调global.Allow返回的permission， 来检查read, write, update, delete, crud分别是否有对应的角色。