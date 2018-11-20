package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	"go-gcs/src/net/http"
	"go-gcs/src/service"
)

// AppRoute will add router
func (a *App) AppRoute() *mux.Router {
	container := restful.NewContainer()
	container.Filter(globalLogging)
	container.Add(newSignedUrlService(a.Service))

	router := mux.NewRouter()
	router.PathPrefix("/v1/").Handler(container)

	return router
}

func newSignedUrlService(sp *service.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/signurl").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Filter(validateTokenMiddleware(sp.Config.JWTSecretKey))
	webService.Route(webService.POST("/").To(http.RESTfulServiceHandler(sp, createGCSSignedUrlHandler)))
	return webService
}
