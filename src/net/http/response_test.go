package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseSuite struct {
	suite.Suite
}

func TestResponseSuite(t *testing.T) {
	suite.Run(t, new(ResponseSuite))
}

func (suite *ResponseSuite) ExampleWriteStatusAndError() {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	if err != nil {
		panic(err)
	}
	WriteStatusAndError(request, recorder, http.StatusBadRequest, errors.New("bad request"))
}

func (suite *ResponseSuite) TestEncodeErrorPayload() {
	testCases := []struct {
		cases       string
		contentType string
		expType     string
		expMessage  string
	}{
		{"json", "text/json", "text/json", `{"error":false,"message":""}`},
		{"xml", "text/xml", "text/xml", `<response><error>false</error><message></message></response>`},
		{"default", "", "application/json", `{"error":false,"message":""}`},
	}

	for _, tc := range testCases {
		errPayload := ErrorPayload{}

		request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
		suite.NoError(err)

		request.Header.Set("Content-Type", tc.contentType)
		out, cType, err := EncodeErrorPayload(request, errPayload)
		suite.Equal(tc.expType, cType)
		suite.NoError(err)
		suite.Equal(tc.expMessage, string(out[:len(out)]))
	}
}

func (suite *ResponseSuite) TestNewErrorPayload() {
	errs := []error{
		fmt.Errorf("Error One"),
		fmt.Errorf("Error Two"),
	}

	err := NewErrorPayload(errs[0], errs[1])
	suite.Equal(errs[0].Error(), err.Message)
	suite.Equal(errs[1].Error(), err.PreviousMessage)
}

func (suite *ResponseSuite) TestWriteStatusAndError() {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
	suite.NoError(err)
	wl, err := WriteStatusAndError(request, recorder, http.StatusForbidden, errors.New("Error one"))
	suite.NoError(err)
	suite.True(wl > 0)

	//Default type is application/json
	suite.Equal(http.StatusForbidden, recorder.Result().StatusCode)
	suite.Equal("application/json", recorder.Header().Get("Content-Type"))
}

func (suite *ResponseSuite) TestSetStatus() {
	testCases := []struct {
		cases      string
		handler    func(req *http.Request, resp http.ResponseWriter, errs ...error) (int, error)
		statusCode int
	}{
		{"Forbidden", Forbidden, http.StatusForbidden},
		{"BadRequest", BadRequest, http.StatusBadRequest},
		{"OK", OK, http.StatusOK},
		{"NotFound", NotFound, http.StatusNotFound},
		{"Unauthorized", Unauthorized, http.StatusUnauthorized},
		{"InternalServerError", InternalServerError, http.StatusInternalServerError},
		{"Conflict", Conflict, http.StatusConflict},
		{"UnprocessableEntity", UnprocessableEntity, http.StatusUnprocessableEntity},
		{"MethodNotAllow", MethodNotAllow, http.StatusMethodNotAllowed},
	}

	for _, tc := range testCases {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest("POST", "http://here.com/v1/signin", nil)
		suite.NoError(err)
		wl, err := tc.handler(request, recorder, errors.New("Failed to do something"))
		suite.NoError(err)
		suite.True(wl > 0)
		suite.Equal(tc.statusCode, recorder.Result().StatusCode)
	}
}
