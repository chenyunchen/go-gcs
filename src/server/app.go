package server

import (
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	"go-gcs/src/config"
	"go-gcs/src/service"
)

// App is the structure to set config & service provider of APP
type App struct {
	Config  config.Config
	Service *service.Container
}

// LoadConfig consumes a string of path to the json config file and read config file into Config.
func (a *App) LoadConfig(configPath string) *App {
	if configPath == "" {
		log.Fatal("-config option is required.")
	}

	a.Config = config.MustRead(configPath)
	return a
}

// Start consumes two strings, host and port, invoke service initilization and serve on desired host:port
func (a *App) Start(host, port string) error {

	a.InitilizeService()

	bind := net.JoinHostPort(host, port)

	return http.ListenAndServe(bind, a.AppRoute())
}

// InitilizeService weavering services with global variables inside server package
func (a *App) InitilizeService() {
	a.Service = service.New(a.Config)
}
