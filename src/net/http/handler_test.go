package http

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go-gcs/src/net/context"
)

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) TestRESTfulServiceHandler() {
	data := 0
	var handler = func(*context.Context) {
		data++
	}

	routeHandler := RESTfulServiceHandler(nil, handler)
	suite.Equal(0, data)
	routeHandler(nil, nil)
	suite.Equal(1, data)
}
