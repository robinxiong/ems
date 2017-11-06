package admin

import (
	"net/http"
	"strings"
	"log"
)
//admin是相关的路由

//Router
type Router struct {
	Prefix string
	routers map[string][]interface{}
}

func (admin *Admin) NewServeMux(prefix string) http.Handler{
	return http.NewServeMux()
}


func (admin *Admin) MountTo(mountTo string, mux *http.ServeMux){
	prefix := "/" + strings.Trim(mountTo, "/")
	serveMux := admin.NewServeMux(prefix)
	log.Print(serveMux)
}