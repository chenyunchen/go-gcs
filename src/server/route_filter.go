package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/emicklei/go-restful"
	"go-gcs/src/logger"
	response "go-gcs/src/net/http"
)

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	logger.Infof("%s %s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func validateTokenMiddleware(jwtSecretKey string) func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		token, err := request.ParseFromRequest(req.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretKey), nil
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
			logger.Infof("Unauthorized access to this resource")
			resp.WriteHeaderAndEntity(http.StatusUnauthorized,
				response.ActionResponse{
					Error:   true,
					Message: "Unauthorized access to this resource",
				})
			return
		}
	}
}
