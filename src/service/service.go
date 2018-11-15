package service

import (
	log "github.com/sirupsen/logrus"
	"go-gcs/src/config"
	"go-gcs/src/service/googlecloudstorageprovider"
	"gopkg.in/go-playground/validator.v9"
)

// Container is the structure for container
type Container struct {
	Config             config.Config
	GoogleCloudStorage *googlecloudstorageprovider.Service
	Validator          *validator.Validate
}

// New will create container
func New(cf config.Config) *Container {
	log.Info("Reading google cloud storage config file: ", cf.GoogleCloudStorage.ServiceAccountPEM)
	googlecloudstorage := googlecloudstorageprovider.New(cf.GoogleCloudStorage.ServiceAccountPEM)

	validate := validator.New()

	return &Container{
		Config:             cf,
		GoogleCloudStorage: googlecloudstorage,
		Validator:          validate,
	}
}

// NewContainer will new a container
func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
