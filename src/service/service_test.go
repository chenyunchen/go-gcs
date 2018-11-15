package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
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
