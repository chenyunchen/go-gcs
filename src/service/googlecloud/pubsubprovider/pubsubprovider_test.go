package pubsubprovider

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/config"
)

type GoogleCloudPubSubProviderSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudPubSubProviderSuite))
}

func (suite *GoogleCloudPubSubProviderSuite) TestNewService() {
	cf := config.MustRead("../../../../config/testing.json")
	service := New(cf.GoogleCloud)
	suite.NotNil(service)
}
