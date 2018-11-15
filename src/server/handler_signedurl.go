package server

import (
	"go-gcs/src/entity"
	"go-gcs/src/googlecloudstorage"
	"go-gcs/src/net/context"
	response "go-gcs/src/net/http"
)

func createGCSSignedUrlHandler(ctx *context.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	r := entity.SignedUrlRequest{}
	if err := req.ReadEntity(&r); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	if err := sp.Validator.Struct(r); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	var signedUrl entity.SignedUrl
	var err error
	switch r.Tag {
	case "single":
		signedUrl, err = googlecloudstorage.CreateGCSSingleSignedUrl(sp, r.FileName, r.ContentType, r.Payload)
	case "group":
		signedUrl, err = googlecloudstorage.CreateGCSGroupSignedUrl(sp, r.FileName, r.ContentType, r.Payload)
	}
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(signedUrl)
}
