package context

import (
	restful "github.com/emicklei/go-restful"
	"go-gcs/src/service"
)

// Context is the struct to combine the restful message with our own serviceProvider
type Context struct {
	ServiceProvider *service.Container
	Request         *restful.Request
	Response        *restful.Response
}
