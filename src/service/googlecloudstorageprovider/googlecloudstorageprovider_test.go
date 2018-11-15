package googlecloudstorageprovider

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GoogleCloudStorageProviderSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudStorageProviderSuite))
}

func (suite *GoogleCloudStorageProviderSuite) TestNewService() {
	service := New("../../key.pem")
	suite.NotNil(service)
}
