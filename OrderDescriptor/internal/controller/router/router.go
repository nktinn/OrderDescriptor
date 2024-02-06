package router

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/service"
)

type Router struct {
	services *service.Services
}

func NewRouter(services *service.Services) *router.Router {
	// Constructor
	r := &Router{
		services: services,
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
