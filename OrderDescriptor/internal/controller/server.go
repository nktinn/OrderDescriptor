package controller

import (
	"github.com/valyala/fasthttp"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/config"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/controller/router"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/service"
)

func StartServer(cfg config.HTTPConfig, services *service.Services) {
	routes := router.NewRouter(services)
	_ = fasthttp.ListenAndServe(cfg.IP+cfg.Port, routes.Handler)
}
