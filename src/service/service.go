package service

import (
	"context"

	"go-gcs/src/config"
	"go-gcs/src/logger"
	"go-gcs/src/service/googlecloud/pubsubprovider"
	"go-gcs/src/service/googlecloud/storageprovider"
	"gopkg.in/go-playground/validator.v9"
)

// Container is the structure for container
type Container struct {
	Config             config.Config
	GoogleCloudPubSub  *pubsubprovider.Service
	GoogleCloudStorage *storageprovider.Service
	Validator          *validator.Validate
}

// New will create container
func New(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Reading google cloud pubsub config file.")

	ctx := context.Background()
	googlecloudpubsub := pubsubprovider.New(ctx, cf.GoogleCloud, cf.PubSub)
	googlecloudstorage := storageprovider.New(ctx, cf.GoogleCloud, cf.Storage)
	validate := validator.New()

	return &Container{
		Config:             cf,
		GoogleCloudPubSub:  googlecloudpubsub,
		GoogleCloudStorage: googlecloudstorage,
		Validator:          validate,
	}
}

// NewContainer will new a container
func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}

// NewForTesting will test container for creating a container
func NewForTesting(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Reading google cloud pubsub config file.")

	ctx := context.Background()
	googlecloudpubsub := pubsubprovider.New(ctx, cf.GoogleCloud, cf.PubSub)
	googlecloudstorage := storageprovider.New(ctx, cf.GoogleCloud, cf.Storage)
	validate := validator.New()

	return &Container{
		Config:             cf,
		GoogleCloudPubSub:  googlecloudpubsub,
		GoogleCloudStorage: googlecloudstorage,
		Validator:          validate,
	}
}
