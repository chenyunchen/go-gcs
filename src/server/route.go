package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	"go-gcs/src/net/http"
	"go-gcs/src/service"
	h "net/http"
)

// AppRoute will add router
func (a *App) AppRoute() *mux.Router {
	container := restful.NewContainer()
	container.Filter(globalLogging)
	container.Add(newStorageService(a.Service))
	container.Add(newPreviewService(a.Service))

	router := mux.NewRouter()
	router.HandleFunc("/", func(w h.ResponseWriter, r *h.Request) {})
	router.PathPrefix("/v1/").Handler(container)

	return router
}

func newStorageService(sp *service.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/storage").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/signurl").Filter(validateTokenMiddleware(sp.Config.JWTSecretKey)).To(http.RESTfulServiceHandler(sp, createGCSSignedUrlHandler)))
	webService.Route(webService.POST("/resize/image").To(http.RESTfulServiceHandler(sp, resizeGCSImageHandler)))
	return webService
}

func newPreviewService(sp *service.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/preview").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(http.RESTfulServiceHandler(sp, previewUrlHandler)))
	return webService
}
