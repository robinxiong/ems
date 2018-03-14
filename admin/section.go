package admin

import (
	"strings"
	"fmt"
	"ems/core/utils"
)

//Section 用来组织forms, 使得填写信息的时候更加的清楚
//    product.EditAttrs(
//      &admin.Section{
//      	Title: "Basic Information",
//      	Rows: [][]string{
//      		{"Name"},
//      		{"Code", "Price"},
//      	}},
//      &admin.Section{
//      	Title: "Organization",
//      	Rows: [][]string{
//      		{"Category", "Collections", "MadeCountry"},
//      	}},
//      "Description",
//      "ColorVariations",
//    }

type Section struct {
	Resource *Resource
	Title    string
	Rows     [][]string
}

// ConvertSectionToStrings convert section to strings
func (res *Resource) ConvertSectionToStrings(sections []*Section) []string {
	var columns []string
	for _, section := range sections {
		for _, row := range section.Rows {
			for _, col := range row {
				columns = append(columns, col)
			}
		}
	}
	return columns
}
// ConvertSectionToMetas convert section to metas
func (res *Resource) ConvertSectionToMetas(sections []*Section) []*Meta {
	var metas []*Meta
	for _, section := range sections {
		for _, row := range section.Rows {
			for _, col := range row {
				meta := res.GetMeta(col)
				if meta != nil {
					metas = append(metas, meta)
				}
			}
		}
	}
	return metas
}

//设置IndexAttr时，调用这个函数
func (res *Resource) setSections(sections *[]*Section, values ...interface{}) {
	//即IndexAttr没有传递任何值时，比如view/index/table.tmpl显示color resource
	if len(values) == 0 {
		if len(*sections) == 0 {
			*sections = res.generateSections(res.allAttrs())
		}
	} else {
		panic("没有实现")
	}
}

//attrs []string, 每一个字段都封装成一个section
func (res *Resource) generateSections(values ...interface{}) []*Section {
	var sections []*Section
	var hasColumns, excludedColumns []string
	// 返传values, 使用得最后一个作为唯一的值
	// e.g. Name, Code, -Name (`-Name` will get first and will skip `Name`)

	for i := len(values) - 1; i >= 0; i-- {
		value := values[i]
		if _, ok := value.(*Section); ok {
			panic("don't implements the section")
		} else if column, ok := value.(string); ok {
			if strings.HasPrefix(column, "-") {
				excludedColumns = append(excludedColumns, column)
			} else if !isContainsColumn(excludedColumns, column) {
				sections = append(sections, &Section{Rows: [][]string{{column}}})
			}
			hasColumns = append(hasColumns, column)
		}else if row, ok := value.([]string); ok {
			for j := len(row) - 1; j >= 0; j-- {
				column = row[j]
				sections = append(sections, &Section{Rows: [][]string{{column}}})
				hasColumns = append(hasColumns, column)
			}
		} else {
			utils.ExitWithMsg(fmt.Sprintf("Resource: attributes should be Section or String, but it is %+v", value))
		}

	}
	sections = reverseSections(sections)
	for _, section := range sections {
		section.Resource = res
	}
	return sections
}
func reverseSections(sections []*Section) []*Section {
	var results []*Section
	for i := 0; i < len(sections); i++ {
		results = append(results, sections[len(sections)-i-1])
	}
	return results
}
//【"-name", "-name2"], name
func isContainsColumn(hasColumns []string, column string) bool {
	for _, col := range hasColumns {
		if strings.TrimLeft(col, "-") == strings.TrimLeft(column, "-") {
			return true
		}
	}
	return false
}
