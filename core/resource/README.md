#Resource
Resource是对一个model struct进行封装，以便提供通用的方法来操作不同的model, 比如CallFindMany, CallFindOne, CallSave等
##Meta
定义resource所包含model， 各字段映射成form中的字段，比如数字，单选框，组合框等, 可以手动指定，也可以通过res.GetMeta(name)将一个数据库字段封装成Meta

```go
user := &struct {
		Name  string
		Name2 *string
		}{}
res := resource.New(user)

//手动指定
meta := &resource.Meta{
		Name:         "Name",
		BaseResource: res,
		}
		
//或者像admin/resource.go那样，通过convertSectionToMetas或者allowedSections, 调用res.GetMeta(name string)自动将数据库的字段转换为Meta
```


##Meta_Value
MetaValue是一个struct, 保存将form，JSON, csv files中的数据，转换为 meta values, 它包含field name, field value 以及它所对应的Meta
如果resource是一个嵌套的resource, MetaValues将包含嵌套的metas
##curd.go 
对model struct的操作
