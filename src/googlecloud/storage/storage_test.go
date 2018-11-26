package storage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/config"
	"go-gcs/src/entity"
	"go-gcs/src/service"
)

type GoogleCloudStorageSuite struct {
	suite.Suite
	sp *service.Container
}

func (suite *GoogleCloudStorageSuite) SetupSuite() {
	cf := config.MustRead("../../../config/testing.json")
	suite.sp = service.New(cf)
}

func TestGoogleCloudStorageSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudStorageSuite))
}

func (suite *GoogleCloudStorageSuite) TestSignURL() {
	path := "test/"
	fileName := "cat.jpg"
	contentType := "image/jpeg"
	method := "PUT"
	expires := time.Now().Add(time.Second * 60)

	url, err := SignURL(suite.sp, path, fileName, contentType, method, expires)
	suite.NoError(err)
	suite.NotEqual("", url)
}

func (suite *GoogleCloudStorageSuite) TestCreateGCSSingleSignedUrl() {
	userId := "myAwesomeId"
	fileName := "cat.jpg"
	contentType := "image/jpeg"
	payload := entity.SinglePayload{
		To: "myAwesomeBuddyId",
	}

	p, err := json.Marshal(payload)
	suite.NotNil(p)
	suite.NoError(err)

	signedUrl, err := CreateGCSSingleSignedUrl(suite.sp, userId, fileName, contentType, string(p))
	suite.NotNil(signedUrl)
	suite.NoError(err)
}

func (suite *GoogleCloudStorageSuite) TestCreateGCSGroupSignedUrl() {
	userId := "myAwesomeId"
	fileName := "cat.jpg"
	contentType := "image/jpeg"
	payload := entity.GroupPayload{
		GroupId: "myAwesomeGroupId",
	}

	p, err := json.Marshal(payload)
	suite.NotNil(p)
	suite.NoError(err)

	signedUrl, err := CreateGCSGroupSignedUrl(suite.sp, userId, fileName, contentType, string(p))
	suite.NotNil(signedUrl)
	suite.NoError(err)
}

func (suite *GoogleCloudStorageSuite) TestResizeGCSImage() {

}
