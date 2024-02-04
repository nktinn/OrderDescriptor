package router

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
)

func (r *Router) Render(ctx *fasthttp.RequestCtx) {
	log.Infoln("Received GET page request. Rendering index.html")
	file, err := os.ReadFile("OrderDescriptor/static/index.html")
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		log.Errorf("Unable to read index.html: %v", err)
		return
	}

	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.Write(file)
	log.Infoln("Rendered index.html")
}

func (r *Router) GetOrder(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("id").(string)
	log.Infoln("Received GET order request /", uid)
	order := r.memrepo.GetFullOrder(uid)
	if order == nil {
		ctx.Error("Order not found", fasthttp.StatusNotFound)
		log.Errorf("Unable to find order with uid %s", uid)
		return
	}
	jsonData, _ := json.Marshal(order)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(jsonData)
	log.Infoln("Returned order with uid", uid)
}

func (r *Router) DeleteOrders(ctx *fasthttp.RequestCtx) {
	log.Infoln("Received DELETE request")

	if err := r.memrepo.DeleteAll(); err != nil {
		ctx.Error(fmt.Sprintf("Unable to delete orders: %v", err), fasthttp.StatusInternalServerError)
		log.Errorf("Unable to delete orders: %v", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	log.Infoln("Deleted all orders")
}
