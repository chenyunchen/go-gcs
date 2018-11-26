package storageprovider

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/service/googlecloud"
)

type Config struct {
	GoogleCloud *googlecloud.Config `json:"googlecloud"`
	Storage     *Storage            `json:"storage"`
}

type GoogleCloudStorageProviderSuite struct {
	suite.Suite
	service *Service
}

func (suite *GoogleCloudStorageProviderSuite) SetupSuite() {
	file, err := os.Open("../../../../config/testing.json")
	suite.NoError(err)
	suite.NotNil(file)
	defer file.Close()

	var cf Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cf)
	suite.NoError(err)

	ctx := context.Background()
	service := New(ctx, cf.GoogleCloud, cf.Storage)
	suite.NotNil(service)
}

func TestStorageProviderSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudStorageProviderSuite))
}
