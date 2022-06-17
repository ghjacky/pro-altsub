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

func (sv *HttpServer) sourceRoutes() {
	srg := sv.Engine.Group(generateRoutePath("sources"))
	srg.POST("", handlerv1.AddSource)
	srg.GET("", handlerv1.FetchSources)
}

func (sv *HttpServer) schemaRoutes() {
	srg := sv.Engine.Group(generateRoutePath("schemas"))
	srg.POST("", handlerv1.AddSchema)
	srg.GET("", handlerv1.FetchSchemas)
}

func (sv *HttpServer) eventRoutes() {
	erg := sv.Engine.Group(generateRoutePath("events"))
	erg.POST("", handlerv1.ReceiveRawEvent)
}

func (sv *HttpServer) ruleRoutes() {
	rrg := sv.Engine.Group(generateRoutePath("rules"))
	rrg.POST("", handlerv1.AddRule)
	rrg.GET("", handlerv1.FetchRules)
}

func (sv *HttpServer) receiverRoutes() {

}

func (sv *HttpServer) subscribeRoutes() {

}

func (sv *HttpServer) maintenanceRoutes() {

}
