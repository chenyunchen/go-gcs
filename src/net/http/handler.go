package http

import (
	"github.com/emicklei/go-restful"
	"go-gcs/src/net/context"
	"go-gcs/src/service"
)

// RESTfulContextHandler is the interface for restfuul handler(restful.Request,restful.Response)
type RESTfulContextHandler func(*context.Context)

// RESTfulServiceHandler is the wrapper to combine the RESTfulContextHandler with our serviceprovider object
func RESTfulServiceHandler(sp *service.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := context.Context{
			ServiceProvider: sp,
			Request:         req,
			Response:        resp,
		}
		handler(&ctx)
	}
}
