package assetfs


type Interface interface {
	PrependPath(path string) error
	RegisterPath(path string) error
	Asset(name string) ([]byte, error)
	Glob(pattern string) (matches []string, err error)
	Compile() error

	NameSpace(nameSpace string) Interface
}