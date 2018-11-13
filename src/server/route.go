package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	"go-gcs/src/service"
)

// Context is the struct to combine the restful message with our own service
type Context struct {
	Service  *service.Container
	Request  *restful.Request
	Response *restful.Response
}

// RESTfulContextHandler is the interface for restfuul handler(restful.Request,restful.Response)
type RESTfulContextHandler func(*Context)

// RESTfulServiceHandler is the wrapper to combine the RESTfulContextHandler with our service object
func RESTfulServiceHandler(sp *service.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := Context{
			Service:  sp,
			Request:  req,
			Response: resp,
		}
		handler(&ctx)
	}
}

// AppRoute will add router
func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()
	container.Add(newUploadService(a.Service))

	router.PathPrefix("/v1/").Handler(container)

	return router
}

func newUploadService(sp *service.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/upload").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(RESTfulServiceHandler(sp, createUploadHandler)))
	return webService
}
