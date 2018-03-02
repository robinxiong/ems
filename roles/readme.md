#roles
https://doc.getqor.com/plugins/roles.html
##角色
全局global role中可以通过register注册它所的角色定义，每一个定义包含一个名称比如，admin或者guest, 然后是根据request和user的确定是否为当前role, 通常匹配数据库中的user.Role。
如果request和user函数，返回true, 则返回角色名称admin或者guest(字符串), 即user拥有这些角色, 然后可以在permission中检查角色是否拥有对资源的权限

##权限Permission主要用于资源的定义，一个资源，允许哪些角色可以访问，可以不可以访问

1.权限模式，有预定义5种read, write, update, delete, crud, 我们也可以创建自己的模式customizer
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