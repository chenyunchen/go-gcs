package googlecloudstorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/config"
	"go-gcs/src/service"
)

type GoogleCloudStorageSuite struct {
	suite.Suite
	sp *service.Container
}

func (suite *GoogleCloudStorageSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	suite.sp = service.New(cf)
}

func TestGoogleCloudStorageSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudStorageSuite))
}

func (suite *GoogleCloudStorageSuite) TestSignURL() {
	path := "test/"
	fileName := "test.txt"
	contentType := "text/plain"
	method := "PUT"
	expires := time.Now().Add(time.Second * 60)

	_, err := SignURL(suite.sp, path, fileName, contentType, method, expires)
	suite.NoError(err)
}
