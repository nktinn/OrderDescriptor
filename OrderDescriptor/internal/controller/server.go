package controller

import (
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/config"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/controller/router"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo"
	"github.com/valyala/fasthttp"
)

func StartServer(cfg config.HTTPConfig, repositories *repo.Repositories, memrepository *repo.MemoryRepo) {
	routes := router.NewRouter(repositories, memrepository)
	_ = fasthttp.ListenAndServe(cfg.IP+cfg.Port, routes.Handler)
}
