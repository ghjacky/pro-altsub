package server

import (
	"altsub/base"
	"altsub/server/handlerv1"
	"path"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
	srg.GET("/:id", handlerv1.GetSource)
	srg.GET("/types/all", handlerv1.FetchSourceTypes)
}

func (sv *HttpServer) schemaRoutes() {
	srg := sv.Engine.Group(generateRoutePath("schemas"))
	srg.POST("", handlerv1.AddSchema)
	srg.GET("", handlerv1.FetchSchemas)
	srg.GET("/:id", handlerv1.GetSchema)
}

func (sv *HttpServer) eventRoutes() {
	erg := sv.Engine.Group(generateRoutePath("events"))
	erg.POST("", handlerv1.ReceiveRawEvent)
}

func (sv *HttpServer) ruleRoutes() {
	rrg := sv.Engine.Group(generateRoutePath("rules"))
	rrg.POST("", handlerv1.AddRule)
	rrg.GET("", handlerv1.FetchRules)
	rrg.GET("/:id", handlerv1.GetRule)
	rrg.DELETE("/:id", handlerv1.DeleteRule)
	rrg.POST("/chain", handlerv1.FetchRuleChain)
	rrg.POST("/assign/:id", handlerv1.Assign)
}

func (sv *HttpServer) receiverRoutes() {
	rrg := sv.Engine.Group(generateRoutePath("receivers"))
	rrg.POST("", handlerv1.AddReceiver)
	rrg.GET("", handlerv1.FetchReceivers)
	rrg.GET("/:id", handlerv1.GetReceiver)
	rrg.DELETE("/:id", handlerv1.DeleteReceiver)
	rrg.POST("/subscribe/:id", handlerv1.Subscribe)
}

func (sv *HttpServer) subscribeRoutes() {
	srg := sv.Engine.Group(generateRoutePath("subscribes"))
	srg.GET("", handlerv1.FetchSubscribe)
}

func (sv *HttpServer) maintenanceRoutes() {

}

func (sv *HttpServer) dutyRoutes() {
	drg := sv.Engine.Group(generateRoutePath("duty"))
	drg.GET("", handlerv1.FetchDuties)
	drg.POST("", handlerv1.AddDuty)
}

func (sv *HttpServer) issueRoutes() {
	irg := sv.Engine.Group(generateRoutePath("issues"))
	irg.GET("", handlerv1.FetchIssueHandlings)
	irg.POST("", handlerv1.AddIssueHandling)
	irg.DELETE("/:id", handlerv1.DeleteIssueHandling)
	irg.GET("/:id", handlerv1.GetIssueHandling)
}

func (sv *HttpServer) staticRouter() {
	sv.Engine.Use(gzip.Gzip(gzip.BestCompression, gzip.WithExcludedPathsRegexs([]string{`^/api/.+$`})))
	sv.Engine.Use(static.Serve("/static", static.LocalFile(base.Config.MainConfig.StaticDir, false)))
	sv.Engine.NoRoute(func(ctx *gin.Context) {
		ctx.File(path.Join(base.Config.MainConfig.StaticDir, "index.html"))
	})
}
