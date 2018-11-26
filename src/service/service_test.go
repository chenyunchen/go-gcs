package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/config"
)

type ServiceSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestContainer() {
	container := NewContainer("../../config/testing.json")
	suite.NotNil(container)
}

func (suite *ServiceSuite) TestNewForTesting() {
	cf := config.MustRead("../../config/testing.json")
	sp := NewForTesting(cf)
	suite.NotNil(sp)
}
