package service

import (
	log "github.com/sirupsen/logrus"
	"go-gcs/src/config"
	"go-gcs/src/service/googlecloud/pubsubprovider"
	"gopkg.in/go-playground/validator.v9"
)

// Container is the structure for container
type Container struct {
	Config            config.Config
	GoogleCloudPubSub *pubsubprovider.Service
	Validator         *validator.Validate
}

// New will create container
func New(cf config.Config) *Container {
	log.Info("Reading google cloud pubsub config file.")
	googlecloudpubsub := pubsubprovider.New(cf.GoogleCloud)

	validate := validator.New()

	return &Container{
		Config:            cf,
		GoogleCloudPubSub: googlecloudpubsub,
		Validator:         validate,
	}
}

// NewContainer will new a container
func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
