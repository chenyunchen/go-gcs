package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"go-gcs/src/logger"
	"go-gcs/src/server"
	"syscall"
)

func main() {
	var (
		configPath string
		host       string
		port       string
	)

	flag.StringVar(&configPath, "config", "config/local.json", "config file path")
	flag.StringVar(&host, "host", "0.0.0.0", "hostname")
	flag.StringVar(&port, "port", "7890", "port")
	flag.Parse()

	a := server.App{}
	go a.LoadConfig(configPath).Start(host, port)

	//for process
	stop := make(chan struct{})

	// Stop all listener by catching interrupt signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		logger.Infof("caught signal: %s", sig.String())

		// TODO check requirement if we need this
		// logger.Infof("deleting google cloud pubsub subscription...")
		// ctx := context.Background()
		// subName := a.Config.PubSub.Subscription
		// sub := a.Service.GoogleCloudPubSub.Client.Subscription(subName)
		// sub.Delete(ctx)

		logger.Infof("all service are stopped successfully")
		close(stop)
	}(sigc)

	<-stop
	os.Exit(0)
}
