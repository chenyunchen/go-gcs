package service

import (
	log "github.com/sirupsen/logrus"
	"go-gcs/src/config"
	"go-gcs/src/service/googlecloudstorageprovider"
)

// Container is the structure for container
type Container struct {
	Config             config.Config
	GoogleCloudStorage *googlecloudstorageprovider.Service
}

// New will create container
func New(cf config.Config) *Container {
	log.Info("Reading google cloud storage config file: ", cf.GoogleCloudStorage.ServiceAccountPEM)
	googlecloudstorage := googlecloudstorageprovider.New(cf.GoogleCloudStorage.ServiceAccountPEM)

	return &Container{
		Config:             cf,
		GoogleCloudStorage: googlecloudstorage,
	}
}
