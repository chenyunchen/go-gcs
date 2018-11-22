package server

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// generateToken is for generating token
func generateToken(userId, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * time.Duration(30)).Unix(),
		"iat": time.Now().Unix(),
		"juu": userId,
	}
	return token.SignedString([]byte(secret))
}

// Assert the response code
func assertResponseCode(t *testing.T, expectedCode int, resp *httptest.ResponseRecorder) {
	t.Helper()
	t.Logf("code:%d", resp.Code)
	if expectedCode != resp.Code {
		t.Errorf("status code %d expected.", expectedCode)
		t.Logf("Response:\n%s", resp.Body.String())
	}
}
