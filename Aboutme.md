依赖第三方
=============================================================================
https://github.com/microcosm-cc/bluemonday 防止跨站攻击，消除html语义
github.com/go-chi/chi 路由
https://github.com/gin-gonic/gin  路由
https://github.com/fatih/color ANSI shell颜色
github.com/manveru/gobdd bdd测试框架
https://github.com/azumads/faker 假数据
github.com/mattn/go-sqlite3 导致编译很慢，所以不要在db中导入split3的支持 或者 go install github.com/mattn/go-sqlite3 
https://github.com/theplant/cldr 常规的翻译包，数字，货币，日历
学习资源
=============================================================================
shell语法 http://tldp.org/LDP/abs/html/comparison-ops.html
shell颜色 https://misc.flogisoft.com/bash/tip_colors_and_formatting
https://blog.golang.org/laws-of-reflection 关于reflect的Addr, addr

学习过程
================================

l10n

publish

sorting

validation

publish2 未完成 
serializable_meta 未完成 

oss 

media 未完成 
  
  如果某个字段带了Media或者oss, 则需要调用注册的回调，来保存或者读取图片，文字信息
  需要先了解以下包
  serializable_meta, 它在保存不固定结构时使用，当结构中涉级到media文件，需要调用media的saveAndCropImage回调
  oss (Object Storage Service)将文件保存到文件系统，FTP, 或者云文件
  


模板设置
==================================
auth登陆页面的模板设置在auth.New auth_themes/clean.New中设置模板所在的路径, auth/controller.go中定义了serveHTTP
