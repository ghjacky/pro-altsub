package server

import "github.com/gin-gonic/gin"

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
	sv.sourceRoutes()
	sv.schemaRoutes()
	sv.eventRoutes()
	sv.receiverRoutes()
	sv.subscribeRoutes()
	sv.maintenanceRoutes()
}

func (sv *HttpServer) RunForever() error {
	return sv.Engine.Run(sv.addr)
}
