package server

import (
	"altsub/server/handlerv1"
	"path"
)

const (
	routePathPrefix = "/api"
)

func generateRoutePath(name string) string {
	return path.Join(routePathPrefix, name)
}

func (sv *HttpServer) schemaRoutes() {

}

func (sv *HttpServer) eventRoutes() {

}

func (sv *HttpServer) receiverRoutes() {

}

func (sv *HttpServer) subscribeRoutes() {

}

func (sv *HttpServer) maintenanceRoutes() {

}

func (sv *HttpServer) sourceRoutes() {
	srg := sv.Engine.Group(generateRoutePath("sources"))
	srg.POST("", handlerv1.AddSource)
}
