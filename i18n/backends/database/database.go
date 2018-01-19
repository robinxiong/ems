package database
type Translation struct {
	Locale string `sql:"size:12;"`
	Key    string `sql:"size:4294967295;"`
	Value  string `sql:"size:4294967295"`
}