package main

import (
	"text/template"

	"fmt"
)

type Render struct {
	*Config
	funcMaps template.FuncMap
}
// Config render config
type Config struct {
	ViewPaths []string
	DefaultLayout string
}

func main() {
	r := Render{
		Config: &Config{
			ViewPaths:[]string{"hello", "world"},
		},
	}
	fmt.Println(len(r.ViewPaths))
	fmt.Println(r.ViewPaths)
}
