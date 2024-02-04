package router

import (
	"github.com/fasthttp/router"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo"
	"github.com/valyala/fasthttp"
)

type Router struct {
	repos   *repo.Repositories
	memrepo *repo.MemoryRepo
}

func NewRouter(repositories *repo.Repositories, memrepository *repo.MemoryRepo) *router.Router {
	// Constructor
	r := &Router{
		repos:   repositories,
		memrepo: memrepository,
	}

	routes := router.New()

	// Render page
	routes.GET("/", r.Render)
	routes.GET("/favicon.ico", func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	})

	// Orders routes
	routes.GET("/{id}", r.GetOrder)
	routes.DELETE("/", r.DeleteOrders)

	return routes
}
