package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	response "go-gcs/src/net/http"
)

var (
	SecretKey = "R5fEQ8weZyC4p8BpknA3UYMt4AU9cA8X"
)

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Info(req.Request.Method, ", ", req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func validateTokenMiddleware(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token, err := request.ParseFromRequest(req.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// save user ID to requests attributes
			req.SetAttribute("userId", claims["juu"])
			chain.ProcessFilter(req, resp)
		} else {
			resp.WriteHeaderAndEntity(http.StatusUnauthorized,
				response.ActionResponse{
					Error:   true,
					Message: "Token is invalid",
				})
			return
		}
	} else {
		log.Info("Unauthorized access to this resource")
		resp.WriteHeaderAndEntity(http.StatusUnauthorized,
			response.ActionResponse{
				Error:   true,
				Message: "Unauthorized access to this resource",
			})
		return
	}
}
