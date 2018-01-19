package resource

type MetaValues struct {
	Values []*MetaValue
}
// MetaValue是一个struct类型，用于保存资源的元数据，这些元数据来自于http表单，json, csv. 它将包含字段名字，字段值以及可配置的元数据, 如果是一个嵌套资源，将会包含子资源的meta信息
type MetaValue struct {
	Name string
	Value interface{}
	Index int //排序
	Meta Metaor

}
