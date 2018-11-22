package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/suite"
	"go-gcs/src/config"
	"go-gcs/src/entity"
	"go-gcs/src/service"
	"go-gcs/src/service/googlecloud/storageprovider"
)

type StorageSuite struct {
	suite.Suite
	sp        *service.Container
	wc        *restful.Container
	JWTBearer string
}

func (suite *StorageSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := service.NewForTesting(cf)

	// init service provider
	suite.sp = sp
	// init restful container
	suite.wc = restful.NewContainer()

	storageService := newStorageService(suite.sp)

	suite.wc.Add(storageService)

	token, err := generateToken("myAwesomeId", sp.Config.JWTSecretKey)
	suite.NotEmpty(token)
	suite.NoError(err)
	suite.JWTBearer = "Bearer " + token
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}

func (suite *StorageSuite) TestCreateGCSSignedUrl() {
	to := "myAwesomeBuddyId"
	singlePayload := entity.SinglePayload{
		To: to,
	}
	bodyBytes, err := json.MarshalIndent(singlePayload, "", "  ")
	suite.NoError(err)

	fileName := "cat.jpg"
	contentType := "image/jpeg"
	tag := "single"
	signedUrlRequest := entity.SignedUrlRequest{
		FileName:    fileName,
		ContentType: contentType,
		Tag:         tag,
		Payload:     string(bodyBytes),
	}
	bodyBytes, err = json.MarshalIndent(signedUrlRequest, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storage/signurl", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Authorization", suite.JWTBearer)
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	signedUrl := entity.SignedUrl{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &signedUrl)
	suite.NoError(err)
	suite.Equal(contentType, signedUrl.UploadHeaders.ContentType)
	suite.Equal(suite.sp.GoogleCloudStorage.Config.ContentLengthRange, signedUrl.UploadHeaders.ContentLengthRange)
}

func (suite *StorageSuite) TestResizeGCSImage() {
	bucket := suite.sp.GoogleCloudStorage.Config.Bucket
	path := "test/cat.jpg"
	filePath := "../../test/image/cat.jpg"
	err := suite.sp.GoogleCloudStorage.Upload(bucket, path, filePath)
	suite.NoError(err)

	gcsPublicBaseUrl := storageprovider.GoogleCloudStoragePublicBaseUrl
	url := fmt.Sprintf("%s/%s/%s", gcsPublicBaseUrl, bucket, path)
	contentType := "image/jpeg"
	resizeImageRequest := entity.ResizeImageRequest{
		Url:         url,
		ContentType: contentType,
	}
	bodyBytes, err := json.MarshalIndent(resizeImageRequest, "", "  ")
	suite.NoError(err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://localhost:7890/v1/storage/resize/image", bodyReader)
	suite.NoError(err)

	httpRequest.Header.Add("Content-Type", "application/json")
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	resizeImage := entity.ResizeImage{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &resizeImage)
	suite.NoError(err)

	imageResizeBucket := suite.sp.GoogleCloudStorage.Config.ImageResizeBucket
	url = fmt.Sprintf("%s/%s/%s", gcsPublicBaseUrl, imageResizeBucket, path)
	suite.Equal(url, resizeImage.Origin)
	suite.Equal(url+"_100", resizeImage.ThumbWidth100)
	suite.Equal(url+"_150", resizeImage.ThumbWidth150)
	suite.Equal(url+"_300", resizeImage.ThumbWidth300)
	suite.Equal(url+"_640", resizeImage.ThumbWidth640)
	suite.Equal(url+"_1080", resizeImage.ThumbWidth1080)

	err = suite.sp.GoogleCloudStorage.Delete(bucket, path)
	suite.NoError(err)
	err = suite.sp.GoogleCloudStorage.Delete(imageResizeBucket, path)
	suite.NoError(err)
	err = suite.sp.GoogleCloudStorage.Delete(imageResizeBucket, path+"_100")
	suite.NoError(err)
	err = suite.sp.GoogleCloudStorage.Delete(imageResizeBucket, path+"_150")
	suite.NoError(err)
	err = suite.sp.GoogleCloudStorage.Delete(imageResizeBucket, path+"_300")
	suite.NoError(err)
	err = suite.sp.GoogleCloudStorage.Delete(imageResizeBucket, path+"_640")
	suite.NoError(err)
	err = suite.sp.GoogleCloudStorage.Delete(imageResizeBucket, path+"_1080")
	suite.NoError(err)
}
