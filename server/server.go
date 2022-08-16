package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Engine *gin.Engine
	addr   string
}

func NewServer(addr string, mode string) *HttpServer {
	gin.SetMode(mode)
	return &HttpServer{
		Engine: gin.New(),
		addr:   addr,
	}
}

func (sv *HttpServer) RegisterRoutes() {
	sv.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	sv.sourceRoutes()
	sv.schemaRoutes()
	sv.eventRoutes()
	sv.ruleRoutes()
	sv.receiverRoutes()
	sv.subscribeRoutes()
	sv.maintenanceRoutes()
	sv.dutyRoutes()
	sv.staticRouter()
}

func (sv *HttpServer) RunForever() error {
	return sv.Engine.Run(sv.addr)
}
