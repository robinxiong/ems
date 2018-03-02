package admin

type ActionArgument struct {
	PrimaryValues []string
	Context *Context
	Argument interface{}
	SkipDefaultResponse bool
}

//action action定义
type Action struct {
	Name string
	Label string
	Method string
}
