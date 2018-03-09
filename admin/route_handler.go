package admin

import (
	"strings"
	"ems/roles"
)

type RouteConfig struct {
	Resource *Resource
	Permissioner   HasPermissioner
	PermissionMode roles.PermissionMode
}

type requestHandler func(c *Context)

type routeHandler struct {
	Path string
	Handle requestHandler
	Config *RouteConfig
}


func newRouteHandler(path string, handle requestHandler, configs ...*RouteConfig) *routeHandler {
	handler := &routeHandler{
		Path:   "/" + strings.TrimPrefix(path, "/"),
		Handle: handle,
	}

	for _, config := range configs {
		handler.Config = config
	}

	if handler.Config == nil {
		handler.Config = &RouteConfig{}
	}

	return handler
}