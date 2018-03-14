package resource

import (
	"ems/core"
	"ems/roles"
)

//CallFindMany
//第一个为resource.NewSlice()返回的值, 它是一个指向到数组的指针，数组又是resource model struct的类型, e.g. &[]*Product
func (res *Resource) CallFindMany(i interface{}, context *core.Context) error {
	return res.FindManyHandler(i, context)
}

func (res *Resource) CallFindOne(i interface{}, value *MetaValues, context *core.Context) error {
	panic("implement me")
}

func (res *Resource) CallSave(i interface{}, context *core.Context) error {
	panic("implement me")
}

func (res *Resource) CallDelete(i interface{}, context *core.Context) error {
	panic("implement me")
}


func (res *Resource) findManyHandler(result interface{}, context *core.Context) error {
	if res.HasPermission(roles.Read, context) {
		//context.GetDB()是返回当前context所局的db, 它在searcher的parseContext中设置了where, order by, limit
		db := context.GetDB()
		//查看db中是否设置了qor:getting_total_count, 可以查看admin.searcher的parseContext
		if _, ok := db.Get("ems:getting_total_count"); ok {
			return context.GetDB().Count(result).Error
		}
		return context.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result).Error
	}
	return roles.ErrPermissionDenied
}