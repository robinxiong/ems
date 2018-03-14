package admin

import (
	"ems/core/resource"
	"ems/core"
	"strconv"
)



// PaginationPageCount default pagination page count
var PaginationPageCount = 20

// Pagination is used to hold pagination related information when rendering tables
type Pagination struct {
	Total       int
	Pages       int
	CurrentPage int
	PerPage     int
}

type Searcher struct {
	*Context
	//admin/products?scopes=Women
	//admin/orders?scopes=processing
	scopes []*Scope //查询的范围，比如产品的
	//admin/product_images?filters[SelectedType].Value=image&filters[Color].Value=2&filters[Category].Value=2
	filters map[*Filter]*resource.MetaValue  //core/resource/meta_value.go
	Pagination Pagination  //parseContext()将调用数据库，获取总的数量
}

//FindMany 基于当前的条件，查找所有的记录
func (s *Searcher) FindMany() (interface{}, error) {
	var (
		err error
		//Search
		context = s.parseContext(true)
		result  = s.Resource.NewSlice()
	)

	if context.HasError() {
		return result, context.Errors
	}

	err = s.Resource.CallFindMany(result, context)
	return result, err
}

func (s *Searcher) clone() *Searcher {
	return &Searcher{Context:s.Context, scopes:s.scopes, filters: s.filters}
}

//解析request中的scopes和filters参数
func (s *Searcher) parseContext(withDefaultScope bool) *core.Context{
	var (
		searcher = s.clone()  //复制searcher
		context = searcher.Context.Context.Clone() //searcher保存的是admin.Context, admin.Context内嵌了core.Context, 调用core.Context.Clone方法
	)

	if context != nil && context.Request != nil {
		//r.ParseMultipartForm()在route.SeverHTTP中解析，获取到多个scopes的值
		//?scopes=1&scopes=2&scopes=3
		scopes := context.Request.Form["scopes"]
		searcher = searcher.Scope(scopes...)
	}
	//todo:filter

	//pagination
	db := context.GetDB()
	//调用resource/crud.go中的findManyHandler, 查询当前filter和scope下的记录的总数
	context.SetDB(db.Model(s.Resource.Value).Set("ems:getting_total_count", true))
	s.Resource.CallFindMany(&s.Pagination.Total, context)


	//获取request中的page, 即当前页
	if s.Pagination.CurrentPage == 0 {
		if s.Context.Request != nil {
			if page, err := strconv.Atoi(s.Context.Request.Form.Get("page")); err == nil {
				s.Pagination.CurrentPage = page
			}
		}

		if s.Pagination.CurrentPage == 0 {
			s.Pagination.CurrentPage = 1
		}
	}

	//是否在request指定了每页的数量，如果没有，使用默认的20
	if s.Pagination.PerPage == 0 {
		if perPage, err := strconv.Atoi(s.Context.Request.Form.Get("per_page")); err == nil {
			s.Pagination.PerPage = perPage
		} else if s.Resource.Config.PageCount > 0 {
			s.Pagination.PerPage = s.Resource.Config.PageCount
		} else {
			s.Pagination.PerPage = PaginationPageCount
		}
	}



	if s.Pagination.CurrentPage > 0 {
		s.Pagination.Pages = (s.Pagination.Total-1)/s.Pagination.PerPage + 1

		db = db.Limit(s.Pagination.PerPage).Offset((s.Pagination.CurrentPage - 1) * s.Pagination.PerPage)
	}

	context.SetDB(db)

	return context
}

//Scope 通过定义的scopes来过滤，names是request.Form["scopes"]的值
func (s *Searcher) Scope(names ...string) *Searcher {
	newSearcher := s.clone()
	//查找资源中是否有定义这个scope，比如查找scopes=Men， 而这个scope.Default为空
	for _, name := range names {
		for _, scope := range s.Resource.scopes {
			if scope.Name == name && !scope.Default {
				//todo: 完成scope绑定到db
			}
		}
	}
	return newSearcher
}