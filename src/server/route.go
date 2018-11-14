package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	"go-gcs/src/net/http"
	"go-gcs/src/service"
)

// AppRoute will add router
func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()
	container.Add(newSignedUrlService(a.Service))

	router.PathPrefix("/v1/").Handler(container)

	return router
}

func newSignedUrlService(sp *service.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/signurl").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(http.RESTfulServiceHandler(sp, createGCSSignedUrlHandler)))
	return webService
}
