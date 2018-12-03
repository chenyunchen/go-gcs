package pubsubprovider

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"go-gcs/src/logger"
	"go-gcs/src/service/googlecloud"
	"go-gcs/src/service/googlecloud/storageprovider"
	"google.golang.org/api/option"
)

// PubSub is the structure for config
type PubSub struct {
	Topic        string `json:"topic"`
	Subscription string `json:"subscription"`
}

// Service is the structure for service
type Service struct {
	Config  *PubSub
	Client  *pubsub.Client
	Context context.Context
}

// GoogleCloudStorageNotification is the structure for notification
type GoogleCloudStorageNotification struct {
	Name        string `json:"name" validate:"required"`
	Bucket      string `json:"bucket" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
}

// NotifyFromGCSStorage will call if google cloud storage object update
func (s *Service) NotifyFromGCSStorage(sp *storageprovider.Service) {
	var sub *pubsub.Subscription
	sub, err := s.Client.CreateSubscription(s.Context, s.Config.Subscription, pubsub.SubscriptionConfig{
		Topic:       s.Client.Topic(s.Config.Topic),
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		logger.Warnf("error while create google cloud pubsub subscription: %s", err)
		logger.Info("try to use the exist subscription...")
		sub = s.Client.Subscription(s.Config.Subscription)
	}

	var mu sync.Mutex
	cctx, cancel := context.WithCancel(s.Context)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		if msg.Attributes["eventType"] == "OBJECT_FINALIZE" {
			gcsNotification := GoogleCloudStorageNotification{}
			json.Unmarshal(msg.Data, &gcsNotification)

			if gcsNotification.ContentType == "image/jpg" ||
				gcsNotification.ContentType == "image/jpeg" ||
				gcsNotification.ContentType == "jpeg" ||
				gcsNotification.ContentType == "image/png" {
				go sp.ResizeMultiImageSizeAndUpload(gcsNotification.ContentType, gcsNotification.Bucket, gcsNotification.Name)
			}
		}
		mu.Lock()
		defer mu.Unlock()
	})
	if err != nil {
		logger.Warnf("error while create google cloud pubsub notify: %s", err)
		cancel()
	}
}

// New will reture a new service
func New(ctx context.Context, googleCloudConfig *googlecloud.Config, pubsubConfig *PubSub) *Service {
	plan, err := json.Marshal(googleCloudConfig)
	if err != nil {
		logger.Warnf("error while read config file: %s", err)
	}

	client, err := pubsub.NewClient(ctx, googleCloudConfig.ProjectId, option.WithCredentialsJSON(plan))
	if err != nil {
		logger.Warnf("error while create google cloud pubsub client: %s", err)
	}

	return &Service{
		Config:  pubsubConfig,
		Client:  client,
		Context: ctx,
	}
}
