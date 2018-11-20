package pubsubprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	log "github.com/sirupsen/logrus"
	"go-gcs/src/service/googlecloud"
	"google.golang.org/api/option"
)

// PubSub is the structure for config
type PubSub struct {
	Topic        string `json:"topic"`
	Subscription string `json:"subscription"`
}

// Service is the structure for service
type Service struct {
	Client *pubsub.Client
}

// Notify will call if google cloud storage object update
func Notify(ctx context.Context, client *pubsub.Client, topic string, subscription string) {
	sub, err := client.CreateSubscription(ctx, subscription, pubsub.SubscriptionConfig{
		Topic:       client.Topic(topic),
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		log.Warn("error while create google cloud pubsub subscription: ", err)
	}

	var mu sync.Mutex
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Println(msg.Attributes["eventType"])
		fmt.Printf("Got message: %q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
	})
	if err != nil {
		log.Warn("error while create google cloud pubsub notify: ", err)
		cancel()
	}
}

// New will reture a new service
func New(googleCloudConfig *googlecloud.Config, pubsubConfig *PubSub) *Service {
	plan, err := json.Marshal(googleCloudConfig)
	if err != nil {
		log.Warn("error while read config file: ", err)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, googleCloudConfig.ProjectId, option.WithCredentialsJSON(plan))
	if err != nil {
		log.Warn("error while create google cloud pubsub client: ", err)
	}

	// go Notify(ctx, client, pubsubConfig.Topic, pubsubConfig.Subscription)

	return &Service{
		Client: client,
	}
}
