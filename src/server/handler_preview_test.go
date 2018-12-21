package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/suite"
	"go-gcs/src/config"
	"go-gcs/src/entity"
	"go-gcs/src/service"
)

type PreviewSuite struct {
	suite.Suite
	wc *restful.Container
	sp *service.Container
}

func (suite *PreviewSuite) SetupSuite() {
	cf := config.MustRead("../../config/testing.json")
	sp := service.NewForTesting(cf)

	// init service provider
	suite.sp = sp
	// init restful container
	suite.wc = restful.NewContainer()

	previewService := newPreviewService(suite.sp)

	suite.wc.Add(previewService)
}

func TestPreviewSuite(t *testing.T) {
	suite.Run(t, new(PreviewSuite))
}

func (suite *PreviewSuite) TestPreviewUrl() {
	httpRequest, err := http.NewRequest("GET", "http://localhost:7890/v1/preview", nil)
	suite.NoError(err)

	url := "https://github.com/chenyunchen"
	query := httpRequest.URL.Query()
	query.Add("url", url)
	httpRequest.URL.RawQuery = query.Encode()
	httpWriter := httptest.NewRecorder()
	suite.wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(suite.T(), http.StatusOK, httpWriter)

	previewedUrl := entity.PreviewedUrl{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &previewedUrl)
	suite.NoError(err)
	suite.Equal("default", previewedUrl.Type)
    // suite.Equal("https://assets-cdn.github.com/favicon.ico", previewedUrl.Icon)
	suite.Equal("GitHub", previewedUrl.Name)
	suite.Equal("chenyunchen - Overview", previewedUrl.Title)
	suite.Equal(url, previewedUrl.Url)
}
