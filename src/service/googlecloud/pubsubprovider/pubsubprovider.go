package pubsubprovider

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	log "github.com/sirupsen/logrus"
	"go-gcs/src/service/googlecloud"
	"google.golang.org/api/option"
)

// Service is the structure for Service
type Service struct {
	PubSubClient *pubsub.Client
}

// New will reture a new service
func New(googleCloudConfig *googlecloud.Config) *Service {
	plan, err := json.Marshal(googleCloudConfig)
	if err != nil {
		log.Warn("error while read config file: ", err)
	}

	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, googleCloudConfig.ProjectId, option.WithCredentialsJSON(plan))
	if err != nil {
		log.Warn("error while create google cloud pubsub client: ", err)
	}

	return &Service{
		PubSubClient: pubsubClient,
	}
}
