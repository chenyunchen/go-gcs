package storageprovider

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/imageresize"
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
	suite.service = service
}

func TestStorageProviderSuite(t *testing.T) {
	suite.Run(t, new(GoogleCloudStorageProviderSuite))
}

func (suite *GoogleCloudStorageProviderSuite) TestUpload() {
	bucket := suite.service.Config.Bucket
	path := "test/cat.jpg"
	filePath := "../../../../test/image/cat.jpg"
	err := suite.service.Upload(bucket, path, filePath)
	suite.NoError(err)
}

func (suite *GoogleCloudStorageProviderSuite) TestUploadImage() {
	contentType := "image/jpeg"
	filePath := "../../../../test/image/cat.jpg"
	img, err := imageresize.ReadImageFile(contentType, filePath)
	suite.NotNil(img)
	suite.NoError(err)

	bucket := suite.service.Config.Bucket
	path := "test/cat.jpg"
	err = suite.service.UploadImage(img, contentType, bucket, path)
	suite.NoError(err)
}

func (suite *GoogleCloudStorageProviderSuite) TestResizeImageAndUpload() {
	contentType := "image/jpeg"
	filePath := "../../../../test/image/cat.jpg"
	img, err := imageresize.ReadImageFile(contentType, filePath)
	suite.NotNil(img)
	suite.NoError(err)

	path := "test/cat.jpg"
	err = suite.service.ResizeImageAndUpload(img, 100, contentType, path)
	suite.NoError(err)
}

func (suite *GoogleCloudStorageProviderSuite) TestResizeMultiImageSizeAndUpload() {
	bucket := suite.service.Config.Bucket
	path := "test/cat.jpg"
	filePath := "../../../../test/image/cat.jpg"
	err := suite.service.Upload(bucket, path, filePath)
	suite.NoError(err)

	contentType := "image/jpeg"
	url, err := suite.service.ResizeMultiImageSizeAndUpload(contentType, bucket, path)
	suite.NoError(err)
	suite.NotEmpty(url)
}
