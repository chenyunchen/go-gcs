package pubsubprovider

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/service/googlecloud"
	"go-gcs/src/service/googlecloud/storageprovider"
)

type Config struct {
	GoogleCloud *googlecloud.Config      `json:"googlecloud"`
	PubSub      *PubSub                  `json:"pubsub"`
	Storage     *storageprovider.Storage `json:"storage"`
}

type GoogleCloudPubSubProviderSuite struct {
	suite.Suite
	service *Service
}

func (suite *GoogleCloudPubSubProviderSuite) SetupSuite() {
	file, err := os.Open("../../../../config/testing.json")
	suite.NoError(err)
	suite.NotNil(file)
	defer file.Close()

	var cf Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cf)
	suite.NoError(err)

	ctx := context.Background()
	service := New(ctx, cf.GoogleCloud, cf.PubSub)
	suite.NotNil(service)
	suite.service = service
}

func (suite *GoogleCloudPubSubProviderSuite) TearDownSuite() {
	subName := suite.service.Config.Subscription
	sub := suite.service.Client.Subscription(subName)
	sub.Delete(suite.service.Context)
}

func TestPubSubProviderSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudPubSubProviderSuite))
}

func (suite *GoogleCloudPubSubProviderSuite) TestNotifyFromGCSStorage() {
	file, err := os.Open("../../../../config/testing.json")
	suite.NoError(err)
	suite.NotNil(file)
	defer file.Close()

	var cf Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cf)
	suite.NoError(err)

	ctx := context.Background()
	storageService := storageprovider.New(ctx, cf.GoogleCloud, cf.Storage)
	suite.NotNil(storageService)

	// Wait for subscription init
	go suite.service.NotifyFromGCSStorage(storageService)
	time.Sleep(5 * time.Second)

	bucket := storageService.Config.Bucket
	path := "test/cat.jpg"
	filePath := "../../../../test/image/cat.jpg"
	err = storageService.Upload(bucket, path, filePath)
	suite.NoError(err)
}
